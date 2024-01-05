package resp

import (
	"bufio"
	"io"

	"github.com/xmopen/golib/pkg/xlogging"
)

// RespType 本次解析类型
type RespType byte

// 解析类型枚举
const (
	// RespTypeUnknown 未知类型
	RespTypeUnknown RespType = '?'
	// RespTypeSimpleString 简单字符串类型
	RespTypeSimpleString RespType = '+'
	// RespTypeError 错误类型
	RespTypeError RespType = '-'
	// RespTypeInteger Integer类型
	RespTypeInteger RespType = ':'
	// ParseTypeString 字符串类型
	RespTypeString RespType = '$'
	// RespTypeArray 数组类型
	RespTypeArray RespType = '*'
)

var (
	// CRLF RESP格式结尾
	CRLF = []byte("\r\n")
)

var (
	parseTable = []Parse{
		parseSimpleString,
		parseError,
		parseInteger,
		parseBulkString,
		parseArray,
	}
)

// Parse 解析函数
type Parse func(xlog *xlogging.Entry, reader *bufio.Reader, pt RespType, body []byte) *Payload

// ParseManager 解析管理器
type ParseManager struct {
	reader *bufio.Reader
}

// New 初始化ParseManager
func New(reader io.Reader) *ParseManager {
	return &ParseManager{
		reader: bufio.NewReader(reader),
	}
}

func (p *ParseManager) DoParse(xlog *xlogging.Entry, body []byte) *Payload {
	// body: xxx\r\n
	//body = bytes.TrimSuffix(body, []byte{'\r', '\n'})
	for _, parse := range parseTable {
		if rsp := parse(xlog, p.reader, RespType(body[0]), body[1:]); rsp != nil {
			return rsp
		}
	}
	// TODO： 待处理
	return nil
}
