package resp

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/xmopen/golib/pkg/xlogging"
)

// IntegerReply integer reply结构
type IntegerReply struct {
	Content int
}

// Bytes 返回IntegerReply对应的字节序
func (i *IntegerReply) Bytes() []byte {
	bytes := []byte(byte(RespTypeInteger))
	return bytes
}

func NewIntegerReply(content int) *IntegerReply {
	return &IntegerReply{
		Content: content,
	}
}

func parseInteger(xlog *xlogging.Entry, reader *bufio.Reader, pt RespType, body []byte) *Payload {
	if pt != RespTypeInteger {
		return nil
	}
	content, err := strconv.Atoi(string(body))
	if err != nil {
		xlog.Errorf("parse integer err:[%+v] source:[%+v]", err, string(body))
		return &Payload{
			Data: NewErrorReply(fmt.Sprintf("illegal number:[%+v]", string(body))),
		}
	}
	return &Payload{
		Data: NewIntegerReply(content),
	}
}
