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
