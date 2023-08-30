// @Title async_server.go
// @Description  基于kinx框架，开发一个tcp服务端
// @Author Zero - 2023/8/12 15:46:09

package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/server"
	"time"
)

// 非阻塞式TCP服务 - 服务端
func main() {
	// 创建服务
	asyncServer := server.NewAsyncServer("[kinx V0.1]", "tcp", "127.0.0.1", 9898)
	// 注册连接初始化事件
	asyncServer.OnPrepare(func(conn netpoll.Connection) context.Context {
		return context.Background()
	})

	// 注册连接完成事件
	asyncServer.OnConnect(func(ctx context.Context, session kiface.ISession) context.Context {
		fmt.Printf("[%s] ID: %d 连接初始化... \n", session.GetRemoteAddr().String(), session.GetSessionID())
		return ctx
	})

	// 注册可读事件处理事件
	asyncServer.OnHandler(func(ctx context.Context, session kiface.ISession) error {
		message, err := session.Read(time.Second)
		if err != nil {
			return err
		}
		fmt.Println("[" + string(message.Payload()) + "]")
		return nil
	})

	// 注册连接关闭事件
	asyncServer.OnClosed(func(conn netpoll.Connection) error {
		fmt.Printf("[%s] 连接关闭... \n", conn.RemoteAddr().String())
		return nil
	})

	// 启动服务
	if err := asyncServer.Run(); err != nil {
		panic(err)
	}
}
