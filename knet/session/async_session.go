// @Title async_session.go
// @Description
// @Author Zero - 2023/8/22 10:35:00

package session

import (
	"context"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/packer"
	"net"
	"time"
)

// AsyncSession 同步非阻塞式会话连接,用于管理客户端的连接，搭配AsyncServer服务端使用;
type AsyncSession struct {
	// 会话ID
	ID uint32
	// 会话的客户端连接，基于netpoll
	Conn netpoll.Connection
	// 会话上下文
	Context context.Context

	// 会话是否开启空闲超时处理
	isIdleTimeout bool
	// 会话空闲超时时间，连接空闲超过该时间强制关闭
	idleTimeout time.Duration
	// 消息封包与解包处理器
	packer kiface.IPacker
}

// NewAsyncSession 创建非阻塞会话AsyncSession
func NewAsyncSession(id uint32, conn netpoll.Connection, ctx context.Context, isIdleTimeout bool, idleTimeout time.Duration) *AsyncSession {
	return &AsyncSession{
		ID:            id,
		Conn:          conn,
		Context:       ctx,
		isIdleTimeout: isIdleTimeout,
		idleTimeout:   idleTimeout,
		packer:        packer.NewAsyncPacker(),
	}
}

// GetConn 获取会话的客户端连接
func (as *AsyncSession) GetConn() net.Conn {
	return as.Conn
}

// GetSessionID 获取会话ID
func (as *AsyncSession) GetSessionID() uint32 {
	return as.ID
}

// GetRemoteAddr 获取会话客户端连接地址
func (as *AsyncSession) GetRemoteAddr() net.Addr {
	return as.Conn.RemoteAddr()
}

// GetContext 获取会话上下文
func (as *AsyncSession) GetContext() context.Context {
	return as.Context
}

// 从客户端连接中读取数据，并且解包(非阻塞式读取)
func (as *AsyncSession) Read(timeout time.Duration) (kiface.IMessage, error) {
	_ = as.Conn.SetReadTimeout(timeout)
	return as.packer.UnPack(as.Conn)
}

// Write 向客户端连接写入数据
func (as *AsyncSession) Write(message kiface.IMessage) error {
	// 消息打包
	pack, _ := as.packer.Pack(message)
	_, err := as.Conn.Writer().WriteBinary(pack)
	return err
}

// Stop 关闭会话连接
func (as *AsyncSession) Stop() {
	// 关闭连接
	_ = as.Conn.Close()
}

// IsClose 当前会话是否已关闭
func (as *AsyncSession) IsClose() bool {
	return as.Conn.IsActive()
}
