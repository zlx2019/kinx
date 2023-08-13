// @Title client.go
// @Description
// @Author Zero - 2023/8/12 17:26:45

package tests

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cloudwego/netpoll"
	"log"
	"os"
	"time"
)

func main() {
	network, addr := "tcp", "127.0.0.1:8989"
	// 1. 连接服务端
	conn, err := netpoll.DialConnection(network, addr, time.Second)
	if err != nil {
		panic("conn server failed: " + err.Error())
	}
	log.Println("conn server successful...")

	// 2.注册事件处理函数
	_ = conn.SetOnRequest(onServerRequest)
	// 注册连接关闭hook函数
	_ = conn.AddCloseCallback(onServerClose)

	// 3. 向服务端写入数据
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		conn.Write([]byte(line))
	}
}

// OnRequest 服务端事件处理函数
func onServerRequest(ctx context.Context, connection netpoll.Connection) error {
	reader := connection.Reader()
	bytes, err := reader.Next(reader.Len())
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

// OnClose 服务端连接关闭hook处理函数
func onServerClose(connection netpoll.Connection) error {
	_ = connection.Close()
	fmt.Println("服务器已关闭连接...")
	return nil
}

