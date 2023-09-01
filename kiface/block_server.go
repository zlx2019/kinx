// @Title block_server.go
// @Description 阻塞服务端抽象接口
// @Author Zero - 2023/8/30 20:41:37

package kiface

import "time"

// IBlockServer 阻塞服务端抽象接口
type IBlockServer interface {
	IServer
	// OnHandler 注册数据处理器
	OnHandler(IHandler)
	// OnOptions 注册服务的配置选项
	OnOptions(options ...BlockServerOption)
	// SetIdleTimeout 开启连接空闲超时，并且设置时长
	SetIdleTimeout(time.Duration)
}
