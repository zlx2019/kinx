// @Title packer.go
// @Description	消息数据包处理器抽象层
// @Author Zero - 2023/8/21 16:00:13

package kiface

import "io"

// IPacker 消息包处理器接口，消息封包与消息拆包
type IPacker interface {
	// Pack 消息打包
	Pack(IMessage) ([]byte, error)
	// UnPack 消息拆包
	UnPack(reader io.Reader) (IMessage, error)
}
