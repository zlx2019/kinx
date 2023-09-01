// @Title message.go
// @Description	 消息数据包实现
// @Author Zero - 2023/8/21 16:36:01

package knet

import "github.com/zlx2019/kinx/kiface"

// Message 消息数据包结构
// 消息序列化结构-> [Len|ID|Payload]
type Message struct {
	// 数据内容长度长度,
	len uint64
	// 消息ID
	id uint64
	// 消息数据内容
	payload []byte
}

// NewMessage 构建一个消息
func NewMessage(id uint64, payload []byte) kiface.IMessage {
	return &Message{
		len:     uint64(len(payload)),
		id:      id,
		payload: payload,
	}
}

func (m *Message) Len() uint64 {
	return m.len
}

func (m *Message) ID() uint64 {
	return m.id
}

func (m *Message) Payload() []byte {
	return m.payload
}

func (m *Message) PutID(id uint64) {
	m.id = id
}

func (m *Message) PutLen(len uint64) {
	m.len = len
}

func (m *Message) PutPayload(payload []byte) {
	m.payload = payload
}
