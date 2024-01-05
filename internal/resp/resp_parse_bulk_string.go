package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/xmopen/golib/pkg/xlogging"
)

// BulkString bulk string struct
type BulkString struct {
	Content []byte
}

// NewBulkString 初始化BulkString结构体
func NewBulkString(content []byte) *BulkString {
	return &BulkString{
		Content: content,
	}
}

// Bytes 返回BulkString对应的字节序
// 格式：$长度/r/n实际内容/r/n
func (b *BulkString) Bytes() []byte {
	bytes := []byte(byte(RespTypeString))
	bytes = append(bytes, []byte(strconv.Itoa(len(b.Content)))...)
	bytes = append(bytes, CRLF...)
	bytes = append(bytes, b.Content...)
	bytes = append(bytes, CRLF...)
	return bytes
}

func parseBulkString(xlog *xlogging.Entry, reader *bufio.Reader, pt RespType, body []byte) *Payload {
	if pt != RespTypeString {
		return nil
	}
	strLen, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil || strLen < -1 {
		xlog.Errorf("parse bulk string err:[%+v] len:[%+v] source:[%+v]", err, strLen, string(body))
		return &Payload{
			Data: NewErrorReply(fmt.Sprintf("illegal bulk string header:[%+v]", string(body))),
		}
	}

	content := make([]byte, strLen)
	readNextLen, err := io.ReadFull(reader, content)
	if err != nil || readNextLen != int(strLen) {
		xlog.Errorf("parse bulk string read next err:[%+v] len:[%+v]", err, strLen)
		return &Payload{
			Data: NewErrorReply(fmt.Sprintf("read next err:[%+v]", err)),
		}
	}
	return &Payload{
		Data: NewBulkString(content),
	}
}
