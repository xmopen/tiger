// Package resp Redis序列化
package resp

// Reply RESP协议消息格式
type Reply interface {
	Bytes() []byte
}

// Payload 消息载体
type Payload struct {
	Data  Reply
	Error error
}

// IsError 判断Payload是否为Error
func (p *Payload) IsError() bool {
	return p.Error != nil
}
