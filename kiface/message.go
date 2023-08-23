// @Title message.go
// @Description  消息数据包抽象层
// @Author Zero - 2023/8/21 16:25:20

package kiface

// IMessage 消息抽象接口
type IMessage interface {
	// Len 获取消息内容的长度 (8 byte) 传输时放在数据头部
	Len() uint64
	// ID 获取消息ID	(8 byte)
	ID() uint64
	// Payload 获取消息内容
	Payload() []byte

	// PutID 设置消息ID
	PutID(uint64)
	// PutLen 设置消息内容的长度
	PutLen(uint64)
	// PutPayload 设置消息的内容
	PutPayload([]byte)
}
