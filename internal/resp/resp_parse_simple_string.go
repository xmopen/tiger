package resp

import (
	"bufio"

	"github.com/xmopen/golib/pkg/xlogging"
)

type StatusReply struct {
	Content string
}

// Bytes 返回字节序列
// +,content,CRLF
func (s *StatusReply) Bytes() []byte {
	bytes := []byte{byte(RespTypeSimpleString)}
	bytes = append(bytes, []byte(s.Content)...)
	bytes = append(bytes, CRLF...)
	return bytes
}

// NewStatusReply 初始化状态返回结构
func NewStatusReply(content string) *StatusReply {
	return &StatusReply{
		Content: content,
	}
}

func parseSimpleString(xlog *xlogging.Entry, reader *bufio.Reader, pt RespType, body []byte) *Payload {
	if pt != RespTypeSimpleString {
		return nil
	}
	return &Payload{
		Data: NewStatusReply(string(body)),
	}
}
