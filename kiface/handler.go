package kiface

import (
	"context"
	"net"
)

// IHandler 服务事件处理器
type IHandler interface {

	// OnConnectHandler 会话连接建立回调方法
	// 服务端与客户端建立连接时回调该函数，返回参数为会话的上下文
	OnConnectHandler(conn net.Conn) context.Context

	// OnClosedHandler 会话连接关闭回调方法
	OnClosedHandler(conn net.Conn) error

	// OnHandler 数据处理方法
	// IHandlerContext 本次数据处理的上下文
	OnHandler(IHandlerContext) error
}

// SuperHandler IHandler的抽象实现，业务处理器继承于此实现后就无需重写所有接口
type SuperHandler struct {
}

func (s *SuperHandler) OnConnectHandler(conn net.Conn) context.Context {
	return context.Background()
}

func (s *SuperHandler) OnClosedHandler(conn net.Conn) error {
	return nil
}

func (s *SuperHandler) OnHandler(ctx IHandlerContext) error {
	return nil
}
