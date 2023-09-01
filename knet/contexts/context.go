// @Title context.go
// @Description
// @Author Zero - 2023/9/1 09:13:13

package contexts

import (
	"context"
	"github.com/zlx2019/kinx/kiface"
)

// HandlerContext 会话处理函数上下文
type HandlerContext struct {
	c context.Context
	// 会话连接信息
	s kiface.ISession
	// 可处理的数据消息
	message kiface.IMessage
}

func (hc *HandlerContext) Put(key, value any) {
	hc.c = context.WithValue(hc.c, key, value)
}

func (hc *HandlerContext) Get(key any) any {
	return hc.c.Value(key)
}

// NewHandlerContext 创建数据处理上下文
func NewHandlerContext(s kiface.ISession, m kiface.IMessage) kiface.IHandlerContext {
	return &HandlerContext{
		s:       s,
		message: m,
		c:       context.Background(),
	}
}

func (hc *HandlerContext) GetConn() kiface.ISession {
	return hc.s
}

func (hc *HandlerContext) GetMessage() kiface.IMessage {
	return hc.message
}
