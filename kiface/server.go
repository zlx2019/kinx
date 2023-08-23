// @Title server.go
// @Description 服务端抽象接口
// @Author Zero - 2023/8/10 19:40:49

package kiface

// IServer Server abstract interface
// 服务端顶级接口
type IServer interface {
	// Run 启动并且运行处理服务
	Run() error
	// Shutdown 关闭服务
	Shutdown() error
}

type IAsyncServer interface {
	// IServer 继承与IServer顶级服务接口
	IServer

	// OnOptions 注册服务端配置
	OnOptions(option ...AsyncServerOption)

	// OnPrepare 注册[连接初始化事件]处理函数
	OnPrepare(OnAsyncPrepareHandler)
	// OnConnect 注册[连接完成事件]处理函数
	OnConnect(OnAsyncConnectHandler)
	// OnReadHandler 注册[连接读取事件]处理函数
	OnReadHandler(OnAsyncReadHandler)
	// OnClosed 注册[连接关闭事件]处理函数
	OnClosed(OnAsyncClosedHandler)
}
