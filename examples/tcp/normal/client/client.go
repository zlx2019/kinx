// @Title client.go
// @Description
// @Author Zero - 2023/8/21 17:43:23

package main

import (
	"bufio"
	"github.com/zlx2019/kinx/knet/message"
	"github.com/zlx2019/kinx/knet/packer"
	"net"
	"os"
	"strings"
)

// 阻塞式TCP服务 - 客户端
func main() {
	// 连接服务端
	conn, _ := net.Dial("tcp", "127.0.0.1:9780")
	// 打包器
	pack := packer.NewNormalPacker()
	stdin := bufio.NewReader(os.Stdin)
	var msgId uint64 = 0
	for {
		line, _, _ := stdin.ReadLine()
		if strings.Contains(string(line), "quit") {
			_ = conn.Close()
			break
		}
		p, _ := pack.Pack(message.NewMessage(msgId, line))
		_, _ = conn.Write(p)
	}
}
