// Package client Tiger Client
package client

import (
	"github.com/xmopen/tiger/internal/command"
	"github.com/xmopen/tiger/internal/protocol/tcontext"
)

// TigerClient Tiger客户端，为每一个Connetion封装一个Client
type TigerClient struct {
	ID      int
	Context *tcontext.Context
	Command *command.CommandManager
}

// New 初始化Tiger客户端
func New(ctx *tcontext.Context) *TigerClient {
	return &TigerClient{
		ID:      generateClientID(),
		Context: ctx,
		Command: command.Manager(),
	}
}

// Exec 客户端执行具体的命令
// 客户端通过CommandManager执行具体的命令
// CommandManager返回执行完具体命令之后的结果(可能是成功的结果,也可能是失败的结果)，由当前客户端写入到Connection中返回给客户端
func (t *TigerClient) Exec(body []byte) {
	reply, err := t.Command.Exec(t.Context.Log(), body)
	if err != nil {
		// 错误处理.
	}
	respBytes := reply.Bytes()
	writeLen, err := t.Context.Connection().Write(respBytes)
	if err != nil {
		// write err 处理
	}
	if writeLen != len(respBytes) {
		// len 处理
	}
	// 命令处理成功
	t.Context.Log().Infof("tiger client exec end,result: success,body:[%+v]", string(body))
}
