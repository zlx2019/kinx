// @Title async_options.go
// @Description
// @Author Zero - 2023/8/23 07:41:07

package async

// AsyncServerOption IAsyncServer非阻塞服务端的配置注册函数
type AsyncServerOption func(server IAsyncServer)

// WithOnAsyncPrepare 设置[客户端连接初始化事件]的处理函数
func WithOnAsyncPrepare(handler OnAsyncPrepareHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnPrepare(handler)
	}
}

// WithOnAsyncConnect 设置[客户端连接完成事件]的处理函数
func WithOnAsyncConnect(handler OnAsyncConnectHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnConnect(handler)
	}
}

// WithOnAsyncHandler 设置服务[连接可读取事件]的处理函数
func WithOnAsyncHandler(handler OnAsyncHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnHandler(handler)
	}
}

// WithOnAsyncClosed 设置[连接关闭事件]的处理函数
func WithOnAsyncClosed(handler OnAsyncClosedHandler) AsyncServerOption {
	return func(server IAsyncServer) {
		server.OnClosed(handler)
	}
}
