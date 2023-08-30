// @Title async_hooks.go
// @Description
// @Author Zero - 2023/8/23 07:24:51

package async

import (
	"context"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
)

// IAsyncServer 非阻塞式服务端的事件处理Hooks函数

// OnAsyncPrepareHandler 非阻塞会话的连接初始化事件
// 返回的context会在 OnAsyncConnectHandler 和 OnAsyncHandler 事件中所传递
type OnAsyncPrepareHandler func(conn netpoll.Connection) context.Context

// OnAsyncConnectHandler 非阻塞会话连接完成事件函数
type OnAsyncConnectHandler func(ctx context.Context, session kiface.ISession) context.Context

// OnAsyncHandler 非阻塞会话读取事件处理函数
type OnAsyncHandler func(ctx context.Context, session kiface.ISession) error

// OnAsyncClosedHandler 非阻塞会话关闭事件
type OnAsyncClosedHandler func(conn netpoll.Connection) error