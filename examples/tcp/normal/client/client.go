// @Title client.go
// @Description
// @Author Zero - 2023/8/21 17:43:23

package main

import (
	"github.com/zlx2019/kinx/knet/message"
	"github.com/zlx2019/kinx/knet/packer"
	"net"
)

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
