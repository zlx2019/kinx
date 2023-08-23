// @Title server.go
// @Description
// @Author Zero - 2023/8/21 14:16:32

package main

import (
	"context"
	"fmt"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/server"
	"net"
)

// NewNormalServer 案例一: 基于事件回调模型开发

func main() {
	// 创建服务
	server := server.NewNormalServer("[kinx V1.0]", "tcp", "127.0.0.1", 9780, MyConnectHandler, MyHandler)
	// 启动服务
	err := server.Run()
	if err != nil {
		panic(err)
	}
}

// MyConnectHandler 连接事件处理函数，初始化会话的上下文
func MyConnectHandler(conn net.Conn) context.Context {
	// 向上下文中放入数据
	return context.WithValue(context.Background(), "SessionName", "zeros")
}

// MyHandler 数据读取事件处理
// session： 触发事件的会话连接
// message:  读取到的消息数据
func MyHandler(session kiface.ISession, message kiface.IMessage) error {
	// 输出消息体
	msg := fmt.Sprintf("ID: %d message: %s", message.ID(), string(message.Payload()))
	fmt.Println(msg)

	// 读取会话上下文中的自定义数据
	name := session.GetContext().Value("SessionName").(string)
	fmt.Println(name)

	// 回写数据
	_ = session.Write(message)
	return nil
}