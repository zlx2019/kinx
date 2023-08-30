// @Title client.go
// @Description
// @Author Zero - 2023/8/30 20:45:26

package main

import (
	"github.com/zlx2019/kinx/knet/message"
	"github.com/zlx2019/kinx/knet/packer"
	"net"
)

// 非阻塞TCP服务 - 客户端
func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:9898")
	packer := packer.NewNormalPacker()

	msg := message.NewMessage(1001, []byte("hello"))
	pack, _ := packer.Pack(msg)
	for i := 0; i < 1; i++ {
		_, _ = conn.Write(pack)
	}
	select {}
}
