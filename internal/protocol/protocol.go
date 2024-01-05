// Package protocol Tiger protocol
package protocol

import (
	"fmt"
	"sync"

	"github.com/xmopen/tiger/internal/protocol/tcpprotocol"
)

// ProtocolType Tiger Protocol类型
type ProtocolType string

// Protocol 类型枚举
const (
	// ProtocolTypeTCP TCP Protocol
	ProtocolTypeTCP ProtocolType = "tcpprotocol"
)

// Protocol Tiger protocol interface
type Protocol interface {
	ListenAndServer(address string) error
	Client() *sync.Map
}

// Handler 具体协议处理请求
type Handler interface {
	Handler()
	OnOpen()
	OnClose()
}

// New 通过protoType初始化Protocol实例
func New(protoType ProtocolType, close chan struct{}) Protocol {
	switch protoType {
	case ProtocolTypeTCP:
		return tcpprotocol.New(close)
	default:
		panic(fmt.Sprintf("no found protocol by prototype:[%+v]\n", protoType))
	}
}
