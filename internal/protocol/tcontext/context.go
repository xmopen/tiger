// Package tcontext 自定义Context
package tcontext

import (
	"context"
	"net"

	"github.com/xmopen/golib/pkg/utils"

	"github.com/xmopen/golib/pkg/xlogging"
)

// Context 自定义上下文
type Context struct {
	ctx  context.Context
	conn net.Conn
}

// New 初始化Tiger Context上下文
func New(ctx context.Context, conn net.Conn) *Context {
	return &Context{
		ctx:  ctx,
		conn: conn,
	}
}

// Log 从上下文中获取日志句柄
func (c *Context) Log() *xlogging.Entry {
	xlog := c.ctx.Value("xlog")
	if xlog != nil {
		return xlog.(*xlogging.Entry)
	}
	xlog = xlogging.Tag(utils.UUID())
	c.ctx = context.WithValue(c.ctx, "xlog", xlog)
	return xlog.(*xlogging.Entry)
}

// Connection 返回当前上下文所持有的connection
func (c *Context) Connection() net.Conn {
	return c.conn
}
