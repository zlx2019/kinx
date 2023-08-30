// @Title async_server.go
// @Description 非阻塞服务端抽象接口
// @Author Zero - 2023/8/30 20:40:07

package async

import "github.com/zlx2019/kinx/kiface"

// IAsyncServer 非阻塞服务端抽象接口
type IAsyncServer interface {
	// IServer 继承与IServer顶级服务接口
	kiface.IServer

	// OnOptions 注册服务端配置
	OnOptions(option ...AsyncServerOption)

	// OnPrepare 注册[连接初始化事件]处理函数
	OnPrepare(OnAsyncPrepareHandler)
	// OnConnect 注册[连接完成事件]处理函数
	OnConnect(OnAsyncConnectHandler)
	// OnHandler 注册[连接读取事件]处理函数
	OnHandler(OnAsyncHandler)
	// OnClosed 注册[连接关闭事件]处理函数
	OnClosed(OnAsyncClosedHandler)
}
