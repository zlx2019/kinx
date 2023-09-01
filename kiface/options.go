// @Title options.go
// @Description
// @Author Zero - 2023/8/30 19:16:46

package kiface

import "time"

// BlockServerOption IBlockServer阻塞式服务端的配置注册函数
type BlockServerOption func(server IBlockServer)

// WithIdleTimeout 设置连接空闲超时时间
func WithIdleTimeout(timeout time.Duration) BlockServerOption {
	return func(server IBlockServer) {
		server.SetIdleTimeout(timeout)
	}
}
