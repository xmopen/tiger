// Package command Tiger具体命令抽象
package command

import (
	"github.com/xmopen/golib/pkg/xlogging"
	"github.com/xmopen/tiger/internal/resp"
)

// Command Tiger命令抽象
type Command interface {
	Exec(xlog *xlogging.Entry, data []byte) resp.Reply
}
