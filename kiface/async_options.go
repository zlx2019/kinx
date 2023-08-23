// @Title async_options.go
// @Description
// @Author Zero - 2023/8/23 07:41:07

package kiface

// AsyncServerOption AsyncServer非阻塞服务端的配置注册函数
type AsyncServerOption func(server IAsyncServer)

// WithOnPrepare 设置[客户端连接初始化事件]的处理函数
func WithOnPrepare(handler OnAsyncPrepareHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnPrepare(handler)
	}
}

// WithOnConnect 设置[客户端连接完成事件]的处理函数
func WithOnConnect(handler OnAsyncConnectHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnConnect(handler)
	}
}

// WithOnReadHandler 设置服务[连接可读取事件]的处理函数
func WithOnReadHandler(handler OnAsyncReadHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnReadHandler(handler)
	}
}

// WithOnClosedHandler 设置[连接关闭事件]的处理函数
func WithOnClosedHandler(handler OnAsyncClosedHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnClosed(handler)
	}
}
