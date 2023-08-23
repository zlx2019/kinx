// @Title hooks.go
// @Description
// @Author Zero - 2023/8/23 08:13:03

package kiface

import (
	"context"
	"net"
)

// 普通服务端的事件处理Hooks函数

// OnConnectHandler 非阻塞会话的连接事件函数
// 服务器与客户端建立连接后，回调该函数, 用于初始化会话的上下文
type OnConnectHandler func(conn net.Conn) context.Context

// OnHandler 会话读取事件函数
// 服务器从会话连接中读取到数据后，回调该执行函数
//
// @param	session		读取的会话连接
// @param	message		读取到的消息体
type OnHandler func(session ISession, message IMessage) error
