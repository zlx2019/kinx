// @Title context.go
// @Description 会话的数据处理上下文接口
// @Author Zero - 2023/9/1 09:01:03

package kiface

// IHandlerContext 会话的一次数据处理上下文接口
type IHandlerContext interface {
	// GetSession 获取会话
	GetSession() ISession
	// GetMessage 获取本次会话接收到的数据
	GetMessage() IMessage
	// Put 设置上下文数据
	Put(key, value any)
	// Get 根据Key获取上下文数据
	Get(any) any
}
