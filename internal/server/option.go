package server

import "github.com/xmopen/tiger/internal/protocol"

// Option Server 配置
type Option func(s *Server)

// WithTCPProtocol 为Server设置Protocol为TCP
func WithTCPProtocol(port int) Option {
	return func(svr *Server) {
		svr.protocol = protocol.ProtocolTypeTCP
		svr.port = port
	}
}

// WithBind 为Server设置Bind属性
func WithBind(bind string) Option {
	return func(svr *Server) {
		svr.bind = bind
	}
}
