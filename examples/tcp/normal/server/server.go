// @Title server.go
// @Description
// @Author Zero - 2023/8/21 14:16:32

package main

import (
	"context"
	"fmt"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet"
	"net"
	"strings"
	"time"
)

// NewNormalServer 案例一: 基于事件回调模型开发
// 阻塞式TCP服务 - 服务端
func main() {
	// 创建服务端
	s := knet.NewNormalServer(
		// 设置协程池数量
		knet.WithPool(1000),
		// 设置连接空闲超时时间
		knet.WithIdleTimeout(time.Hour*30),
		// 设置处理器
		knet.WithHandler(&CustomHandler{}))
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

// OnConnectHandler 会话连接事件
func (c *CustomHandler) OnConnectHandler(conn net.Conn) context.Context {
	fmt.Printf("[%s] 已连接... \n", conn.RemoteAddr())
	return context.Background()
}

// OnHandler 会话有数据处理事件
func (c *CustomHandler) OnHandler(ctx kiface.IHandlerContext) error {
	message := ctx.GetMessage()
	msg := fmt.Sprintf("[%s]: %s", ctx.GetSession().GetRemoteAddr(), string(message.Payload()))
	fmt.Println(msg)

	// 写回消息
	_ = ctx.GetSession().Write(message)

	// 关闭会话
	if strings.Contains(msg, "stop") {
		ctx.GetSession().Stop()
	}
	return nil
}

// OnClosedHandler 会话关闭事件
func (c *CustomHandler) OnClosedHandler(conn net.Conn) error {
	fmt.Printf("[%s] 已关闭... \n", conn.RemoteAddr())
	return nil
}
