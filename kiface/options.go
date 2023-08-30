// @Title options.go
// @Description
// @Author Zero - 2023/8/30 19:16:46

package kiface

import "time"

// BlockServerOption IBlockServer阻塞式服务端的配置注册函数
type BlockServerOption func(server IBlockServer)

// WithOnConnects 设置[连接完成事件]处理函数
func WithOnConnects(handler OnConnectHandler) BlockServerOption {
	return func(server IBlockServer) {
		server.OnConnect(handler)
	}
}

// WithOnHandler 注册[连接可读取事件]处理函数
func WithOnHandler(handler OnHandler) BlockServerOption {
	return func(server IBlockServer) {
		server.OnHandler(handler)
	}
}

// WithOnClosed 注册[连接关闭事件]处理函数
func WithOnClosed(handler OnClosedHandler) BlockServerOption {
	return func(server IBlockServer) {
		server.OnClosed(handler)
	}
}

// WithIdleTimeout 设置连接空闲超时时间
func WithIdleTimeout(timeout time.Duration) BlockServerOption {
	return func(server IBlockServer) {
		server.SetIdleTimeout(timeout)
	}
}
