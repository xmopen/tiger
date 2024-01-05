package tcpprotocol

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xmopen/tiger/internal/client"

	"github.com/xmopen/tiger/internal/protocol/tcontext"

	"github.com/xmopen/golib/pkg/xgoroutine"
	"github.com/xmopen/golib/pkg/xlogging"
)

const (
	tcpProtoStatusUnknown   uint32 = 0x00
	tcpProtoStatusRunning   uint32 = 0x01
	tcpProtoStatusClosed    uint32 = 0x02
	tcpProtoStatusInterrupt uint32 = 0x04
)

// TCPProtocol Tiger tcp protocol
type TCPProtocol struct {
	close              chan struct{}
	status             uint32
	acceptDelay        time.Duration
	acceptDelayMax     time.Duration
	openTCPKeepAlive   bool
	tcpKeepAlivePeriod time.Duration
	openTCPDeadLine    bool
	tcpDeadLine        time.Duration
	allClients         *sync.Map
	xlog               *xlogging.Entry
}

// New 初始化TCP协议
func New(close chan struct{}) *TCPProtocol {
	return &TCPProtocol{
		close:          close,
		status:         tcpProtoStatusUnknown,
		acceptDelay:    0,
		acceptDelayMax: 1 * time.Second,
		allClients:     &sync.Map{},
		xlog:           xlogging.Tag("tiger.tcpprotocol"),
	}
}

// ListenAndServer 监听Port并且运行Server
func (t *TCPProtocol) ListenAndServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	xgoroutine.SafeGoroutine(func() {
		<-t.close
		t.xlog.Warnf("close success")
		// listener.Accept() 将会返回immediately error
		_ = listener.Close()
	})
	wg := &sync.WaitGroup{}
	defer func() {
		t.xlog.Infof("tcp protocol is exiting")
		wg.Wait()
	}()
	t.setStatus(tcpProtoStatusRunning)
	t.xlog.Infof("tiger running success,start listen client: ...")
	for {
		if !t.isRunning() {
			t.xlog.Warnf("tcp protocol is closed,not go accept")
			return nil
		}
		rwc, err := listener.Accept()
		if err != nil {
			if t.isTCPListenerAcceptTimeout(err) {
				t.xlog.Errorf("listener accept timeout err")
				t.sleepBecauseAcceptTimeout()
				continue
			}
			return err
		}
		if tcpconn, ok := rwc.(*net.TCPConn); ok {
			if err := t.setTCPConnetionAttr(tcpconn); err != nil {
				t.xlog.Errorf("set tcp connetion attr err:[%+v],remote ip:[%+v]", err, tcpconn.RemoteAddr())
				return err
			}
		}
		tigerClient := client.New(tcontext.New(context.TODO(), rwc))
		wg.Add(1)
		xgoroutine.SafeGoroutine(func() {
			defer wg.Done()
			newTCPHandler(t, tigerClient).Handler()
		})
	}
}

func (t *TCPProtocol) isRunning() bool {
	return t.status|tcpProtoStatusRunning > 0
}

func (t *TCPProtocol) isTCPListenerAcceptTimeout(err error) bool {
	var netError net.Error
	return errors.As(err, &netError) && netError.Timeout()
}

func (t *TCPProtocol) sleepBecauseAcceptTimeout() {
	time.Sleep(t.nowDelayDuration())
}

func (t *TCPProtocol) nowDelayDuration() time.Duration {
	if t.acceptDelay == t.acceptDelayMax {
		return t.acceptDelay
	}
	if t.acceptDelay == 0 {
		t.acceptDelay = 5 * time.Millisecond
	} else {
		t.acceptDelay *= 2
	}
	if t.acceptDelay > t.acceptDelayMax {
		t.acceptDelay = t.acceptDelayMax
	}
	return t.acceptDelay
}

func (t *TCPProtocol) setTCPConnetionAttr(tcpconn *net.TCPConn) error {
	if t.openTCPKeepAlive {
		if err := tcpconn.SetKeepAlive(true); err != nil {
			return err
		}
		if err := tcpconn.SetKeepAlivePeriod(t.tcpKeepAlivePeriod); err != nil {
			return err
		}
	}
	if t.openTCPDeadLine {
		return tcpconn.SetDeadline(time.Now().Add(t.tcpDeadLine))
	}
	return nil
}

func (t *TCPProtocol) setStatus(status uint32) {
	if (status | tcpProtoStatusRunning | tcpProtoStatusClosed | tcpProtoStatusInterrupt) < 0 {
		panic(fmt.Sprintf("not found status:[%+v]\n", status))
	}
	atomic.StoreUint32(&t.status, status)
}

// Client 返回ClientMap句柄
func (t *TCPProtocol) Client() *sync.Map {
	return t.allClients
}
