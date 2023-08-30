// @Title block_server.go
// @Description 阻塞服务端抽象接口
// @Author Zero - 2023/8/30 20:41:37

package kiface

import "time"

// IBlockServer 阻塞服务端抽象接口
type IBlockServer interface {
	IServer
	// OnOptions 注册服务的配置选项
	OnOptions(options ...BlockServerOption)
	// OnConnect 注册[连接完成事件]处理函数
	OnConnect(OnConnectHandler)
	// OnHandler 注册[连接读取事件]处理函数
	OnHandler(OnHandler)
	// OnClosed 注册[连接关闭事件]处理函数
	OnClosed(OnClosedHandler)
	// SetIdleTimeout 开启连接空闲超时，并且设置时长
	SetIdleTimeout(time.Duration)
}
