package command

import (
	"sync"

	"github.com/xmopen/golib/pkg/xlogging"
	"github.com/xmopen/tiger/internal/resp"
)

var (
	commandStringInstance *commandString
	commandStringOnce     sync.Once
)

type commandString struct {
}

func initCommandString() *commandString {
	commandStringOnce.Do(func() {
		commandStringInstance = &commandString{}
	})
	return commandStringInstance
}

func (c *commandString) Exec(xlog *xlogging.Entry, data []byte) resp.Reply {
	panic("implement me")
}
