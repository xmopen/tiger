package command

import (
	"sync"

	"github.com/xmopen/golib/pkg/xlogging"
	"github.com/xmopen/tiger/internal/resp"
)

var (
	commandManagerInstance *CommandManager
	commandManagerOnce     sync.Once
)

// CommandManager command manager
type CommandManager struct {
	commandTable map[string]Command
}

// Manager 初始化CommandManager并且返回单例
func Manager() *CommandManager {
	commandManagerOnce.Do(func() {
		commandManagerInstance = &CommandManager{
			commandTable: initCommandTable(),
		}
	})
	return commandManagerInstance
}

func initCommandTable() map[string]Command {
	return map[string]Command{
		"1": initCommandString(),
	}
}

// Exec 执行具体的命令,命令由解析出来的body确认
func (c *CommandManager) Exec(xlog *xlogging.Entry, body []byte) (resp.Reply, error) {
	cmd, ok := c.commandTable[""]
	if !ok {
		// 返回命令不存在
		return nil, nil
	}
	return cmd.Exec(xlog, body), nil
}
