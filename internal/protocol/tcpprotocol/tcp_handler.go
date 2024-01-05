package tcpprotocol

import (
	"sync"

	"github.com/xmopen/tiger/internal/resp"

	"github.com/xmopen/golib/pkg/xgoroutine"

	"github.com/xmopen/tiger/internal/client"
)

// TCPHandler tcp handler
type TCPHandler struct {
	tcp             *TCPProtocol
	client          *client.TigerClient
	streamBuffer    chan *resp.Payload
	streamParsePool *sync.Pool
}

func newTCPHandler(tcp *TCPProtocol, client *client.TigerClient) *TCPHandler {
	return &TCPHandler{
		tcp:          tcp,
		client:       client,
		streamBuffer: make(chan *resp.Payload),
		streamParsePool: &sync.Pool{
			New: func() any {
				return newTCPStreamParse()
			},
		},
	}
}

// Handler 处理每一个TCPConnetion信息
// Handler函数不可以异步执行,上层已经是异步了
func (t *TCPHandler) Handler() {
	t.OnOpen()
	stream, _ := t.streamParsePool.Get().(*TCPStreamParse)
	stream.buildWithReader(t.client.Context.Connection()).buildWithBuffer(t.streamBuffer).buildWithDelim('\n').
		buildParse()
	defer func() {
		t.OnClose()
		t.streamParsePool.Put(stream)
	}()
	// 读取客户端的数据
	xgoroutine.SafeGoroutine(func() {
		stream.ParseWithLine(t.client.Context.Log())
	})
	for payload := range t.streamBuffer {
		if payload.IsError() {
			t.client.Context.Log().Errorf("receive payload err:[%+v]", payload.Error)
			continue
		}
		// 客户端执行具体的命令
		t.client.Exec(payload.Data.Bytes())
	}
}

// OnOpen 处理刚建立完成的TCPConnetion
func (t *TCPHandler) OnOpen() {
	t.tcp.allClients.Store(t.client.ID, t.client)
	t.client.Context.Log().Infof("")
}

// OnClose TCPHandler结束时的回调
func (t *TCPHandler) OnClose() {
	t.tcp.allClients.Delete(t.client.ID)
	t.client.Context.Log().Infof("")
}
