// @Title session.go
// @Description	连接会话抽象层
// @Author Zero - 2023/8/21 13:07:16

package kiface

import (
	"context"
	"net"
	"time"
)

// ISession 会话接口
// 将连接抽象为会话，由会话管理连接
type ISession interface {
	// GetConn 获取会话的连接
	GetConn() net.Conn
	// GetSessionID  获取会话ID
	GetSessionID() uint32
	// GetRemoteAddr 获取连接的地址信息
	GetRemoteAddr() net.Addr
	// GetContext 获取会话的上下文
	GetContext() context.Context
	// Read 从连接中读取数据，并且解包为IMessage
	Read(duration time.Duration) (IMessage, error)
	// Write 向连接写入数据包
	Write(message IMessage) error
	// Stop 关闭会话连接
	Stop()
	// IsClose 会话是否已关闭
	IsClose() bool
}
