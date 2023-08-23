// @Title async_server_event.go
// @Description
// @Author Zero - 2023/8/22 19:05:54

package server

import (
	"context"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/session"
	"sync/atomic"
	"time"
)

const (
	// 连接心跳通道在上下文中的Key
	idleTimeoutChKey = "IDLE_TIMEOUT_CH"
	// SessionContextKey 会话在上下文中存储的KEY
	sessionContextKey = "SESSION_CTX_KEY"
)

// 连接开始初始化事件
// 为连接初始化上下文，直接回调注册的 OnPrepare 处理函数
func (as *AsyncServer) onPrepareEvent(conn netpoll.Connection) context.Context {
	return as.onPrepare(conn)
}

// 连接完成事件
// 为连接初始化会话，放入上下文传递，并且回调注册的 onConnect 处理函数
func (as *AsyncServer) onConnectEvent(ctx context.Context, conn netpoll.Connection) context.Context {
	// 创建会话
	s := session.NewAsyncSession(atomic.AddUint32(&as.nextSessionID, 1), conn, ctx, as.isIdleTimeout, as.idleTimeout)

	// 将会话放入连接上下文
	ctx = context.WithValue(ctx, sessionContextKey, s)

	// 注册连接关闭事件函数
	_ = conn.AddCloseCallback(as.onClosedEvent)

	// 是否开启超时处理
	if as.isIdleTimeout {
		ctx = as.timeoutHandler(ctx, s)
	}

	// 回调自定义连接事件函数
	return as.onConnect(ctx, s)
}

// OnRequestEvent 连接可读事件
func (as *AsyncServer) onReadEvent(ctx context.Context, connection netpoll.Connection) error {
	s := ctx.Value(sessionContextKey).(kiface.ISession)
	return as.onRead(ctx, s)
}

// onClosedEvent 连接关闭事件
func (as *AsyncServer) onClosedEvent(connection netpoll.Connection) error {
	return as.onClosed(connection)
}

// 连接空闲超时处理，超时后主动关闭连接，强制下线
func (as *AsyncServer) timeoutHandler(ctx context.Context, s kiface.ISession) context.Context {
	// 连接空闲超时处理，由于目前netpoll不支持连接超时，所以需要手动实现
	ch := make(chan struct{})
	// TODO 后续使用协程池
	go func() {
		for {
			select {
			case <-ch:
			case <-time.After(as.idleTimeout):
				// 超时.关闭连接
				s.Stop()
				close(ch)
			}
		}
	}()
	// 将心跳通道，在上下文中传递
	return context.WithValue(ctx, idleTimeoutChKey, ch)
}
