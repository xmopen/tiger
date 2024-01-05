package tcpprotocol

import (
	"bufio"
	"io"

	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/tiger/internal/resp"
)

// TCPStreamParse tcp解析流
type TCPStreamParse struct {
	delim  byte
	reader *bufio.Reader
	buffer chan *resp.Payload
	parse  *resp.ParseManager
}

func newTCPStreamParse() *TCPStreamParse {
	return &TCPStreamParse{}
}

func (t *TCPStreamParse) buildWithReader(reader io.Reader) *TCPStreamParse {
	t.reader = bufio.NewReader(reader)
	return t
}

func (t *TCPStreamParse) buildWithBuffer(buffer chan *resp.Payload) *TCPStreamParse {
	t.buffer = buffer
	return t
}

func (t *TCPStreamParse) buildWithDelim(delim byte) *TCPStreamParse {
	t.delim = delim
	return t
}

func (t *TCPStreamParse) buildParse() *TCPStreamParse {
	t.parse = resp.New(t.reader)
	return t
}

// ParseWithLine 按照行规则RESP协议进行解析
// RESP定义了五中格式：
// 1、简单字符串(Simple String): 服务器用来返回简单的结果，比如OK等，非二进制安全(二进制内允许出现任何字符\r\n也是可以出现的)，不允许换行
// 2、错误信息(Error)：服务器用来返回简单的错误信息，非二进制安全，不允许换行
// 3、整数(Integer)：Llen、scard等命令的返回值
// 4、字符串(Bulk String)：二进制安全字符串
// 5、数组(Array，Multi Bulk Strings)：客户端发送指令以及Lrange等命令的响应格式
// RESP具体格式：
// 参考：https://redis.io/docs/reference/protocol-spec/
// 参考：https://www.cnblogs.com/Finley/p/11923168.html
// 格式：
// 1、简单字符：以+开始，参考：+OK\r\n
// 2、错误：以-开始，参考：-ERR Invalid Synatx\r\n
// 3、整数: 以:开始，参考：:1\r\n
// 4、字符串：以$开始
// 5、数组：以*开始
func (t *TCPStreamParse) ParseWithLine(xlog *xlogging.Entry) {
	if t.reader == nil || t.buffer == nil {
		panic("tcp stream parse reader or buffer is nil")
	}
	for {
		data, err := t.reader.ReadBytes(t.delim)
		if err != nil {
			t.buffer <- &resp.Payload{
				Error: err,
			}
			return
		}
		if len(data) <= 2 || data[len(data)-2] != '\r' {
			// TODO: 当不满足定义好的规则时，必须告知客户端当前传输的数据有问题
			xlog.Errorf("recv err data:[%+v]", string(data))
			continue
		}
		t.buffer <- t.parse.DoParse(xlog, data)
	}
}
