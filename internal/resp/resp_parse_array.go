package resp

import (
	"bufio"
	"bytes"
	"strconv"

	"github.com/xmopen/golib/pkg/xlogging"
)

// MultiBulkReply ..
type MultiBulkReply struct {
	Content [][]byte
}

// Bytes 构造MultiBulkReply格式
func (m *MultiBulkReply) Bytes() []byte {
	var buffer bytes.Buffer
	buffer.WriteByte(byte(RespTypeArray))
	buffer.WriteString(strconv.Itoa(len(m.Content)))
	buffer.Write(CRLF)
	for _, arg := range m.Content {
		if arg == nil {
			buffer.WriteByte(byte(RespTypeString))
			buffer.WriteString("-1")
		} else {
			buffer.WriteByte(byte(RespTypeString))
			buffer.WriteString(strconv.Itoa(len(arg)))
			buffer.Write(CRLF)
			buffer.WriteString(string(arg))
		}
		buffer.Write(CRLF)
	}
	return buffer.Bytes()
}

// NewMultiBulkReply 初始化MultiBulkReply实例
func NewMultiBulkReply(content [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Content: content,
	}
}

func parseArray(xlog *xlogging.Entry, reader *bufio.Reader, pt RespType, body []byte) *Payload {
	////// 不想学习...... xxxxxxxxxxxxxxxxxxxxxxxxxxxx
	return nil
}
