// Package server tiger server
package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xmopen/golib/pkg/xgoroutine"

	"github.com/xmopen/golib/pkg/xlogging"
	"github.com/xmopen/tiger/internal/protocol"
)

var (
	defaultTigerServerPort     = 8848
	defaultTigerMaxClients     = 1024
	defaultTigerDataBases      = 16
	defaultTigerDir            = ""
	defaultTigerRDBFileName    = ""
	defaultTigerWriteTimeout   = 3 * time.Second
	defaultTigerServerProtocol = protocol.ProtocolTypeTCP
)

// 1、Server通过Protocol接收客户端请求
// 2、protocol按照指定协议将客户端封装好，转给Server
// 3、Server持有客户端连接句柄进行命令操作

// Server Tiger server
type Server struct {
	port         int
	maxClients   int
	dataBases    int
	runID        string
	bind         string
	dir          string
	rdbFileName  string
	close        chan struct{}
	protocol     protocol.ProtocolType
	writeTimeout time.Duration
	xlog         *xlogging.Entry
}

// New 初始化 Tiger Server
func New(options ...Option) *Server {
	svr := &Server{
		close: make(chan struct{}),
		xlog:  xlogging.Tag("tiger.server"),
	}
	for _, opt := range options {
		opt(svr)
	}
	svr.initServerConfig()
	return svr
}

// ListenAndServerWithSignal Tiger监听端口并且运行
func (s *Server) ListenAndServerWithSignal() error {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	xgoroutine.SafeGoroutine(func() {
		sgn := <-sign
		s.Stop()
		s.xlog.Warnf("signal notify sign:[%+v]", sgn)
	})
	return protocol.New(s.protocol, s.close).ListenAndServer(fmt.Sprintf(":%d", s.port))
}

func (s *Server) Stop() {
	close(s.close)
	s.xlog.Infof("close success")
}

// initServerConfig 校验Server必要参数是否为空
func (s *Server) initServerConfig() {
	if s.port == 0 {
		s.port = defaultTigerServerPort
	}
	if s.maxClients == 0 {
		s.maxClients = defaultTigerMaxClients
	}
	if s.dir == "" {
		s.dir = defaultTigerDir
	}
	if s.dataBases == 0 {
		s.dataBases = defaultTigerDataBases
	}
	if s.rdbFileName == "" {
		s.rdbFileName = defaultTigerRDBFileName
	}
	if s.protocol == "" {
		s.protocol = defaultTigerServerProtocol
	}
	if s.writeTimeout == 0 {
		s.writeTimeout = defaultTigerWriteTimeout
	}
}
