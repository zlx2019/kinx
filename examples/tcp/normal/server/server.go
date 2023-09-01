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
	"time"
)

// NewNormalServer 案例一: 基于事件回调模型开发
// 阻塞式TCP服务 - 服务端

func main() {
	// 创建服务端
	s := server.NewNormalServer("[kinx V1.0]", "127.0.0.1", 9780)
	// 设置处理器
	s.OnHandler(&CustomHandler{})
	s.SetIdleTimeout(time.Second * 30)
	// 启动服务
	err := s.Run()
	if err != nil {
		panic(err)
	}
}

// CustomHandler 自定义处理器，要求实现 IHandler 接口
type CustomHandler struct {
	kiface.SuperHandler
}

func (c *CustomHandler) OnConnectHandler(conn net.Conn) context.Context {
	fmt.Printf("[%s] 已连接... \n", conn.RemoteAddr())
	return context.Background()
}

func (c *CustomHandler) OnHandler(ctx kiface.IHandlerContext) error {
	message := ctx.GetMessage()
	msg := fmt.Sprintf("ID: %d message: %s", message.ID(), string(message.Payload()))
	fmt.Println(msg)
	return nil
}

func (c *CustomHandler) OnClosedHandler(conn net.Conn) error {
	fmt.Printf("[%s] 已关闭... \n", conn.RemoteAddr())
	return nil
}
