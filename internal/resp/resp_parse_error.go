package resp

import (
	"bufio"

	"github.com/xmopen/golib/pkg/xlogging"
)

// ErrorReply 错误返回结构
type ErrorReply struct {
	Content string
}

// Bytes 返回ErrorReply Bytes
func (e *ErrorReply) Bytes() []byte {
	bytes := []byte{byte(RespTypeError)}
	bytes = append(bytes, []byte(e.Content)...)
	bytes = append(bytes, CRLF...)
	return bytes
}

// Error 返回Error
func (e *ErrorReply) Error() string {
	return e.Content
}

// NewErrorReply 初始化ErrorReply结构体
func NewErrorReply(content string) *ErrorReply {
	return &ErrorReply{
		Content: content,
	}
}

func parseError(xlog *xlogging.Entry, reader *bufio.Reader, pt RespType, body []byte) *Payload {
	if pt != RespTypeError {
		return nil
	}
	return &Payload{
		Data: NewErrorReply(string(body)),
	}
}
