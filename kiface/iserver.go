// @Title iserver.go
// @Description 服务端抽象接口
// @Author Zero - 2023/8/10 19:40:49

package kiface

// IServer Server abstract interface
// 服务端抽象层 --> 顶级接口
type IServer interface {
	// Run 运行服务
	Run() error
	// Shutdown 关闭服务
	Shutdown()error
}
