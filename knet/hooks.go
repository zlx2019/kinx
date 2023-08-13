// @Title hooks.go
// @Description	连接的生命周期hook处理函数
// @Author Zero - 2023/8/12 20:27:05

package knet

import (
	"context"
	"github.com/cloudwego/netpoll"
	"log"
	"time"
)

// 连接心跳通道在上下文中的Key
const idleTimeoutChKey = "idleTimeoutCh"

// OnPrepare 连接初始化时会执行的处理函数
// 用于在连接初始化时注入自定义准备，这是可选的,但在某些情况下很重要。
// 返回的context上下文将成为 OnConnect和OnRequest的参数。
func (ks *KServer) OnPrepare(connection netpoll.Connection) context.Context {
	// 可以在这里向下传递一些K-V参数
	return nil
}

// OnConnect 连接创建完成后会执行的处理函数
func (ks *KServer) OnConnect(ctx context.Context, connection netpoll.Connection) context.Context {
	log.Printf("[%s] 连接建立成功.\n", connection.RemoteAddr().String())
	// 注册连接关闭处理函数
	_ = connection.AddCloseCallback(ks.OnClose)

	// 连接空闲超时处理
	if ks.enableIdleTimeout {
		ctx = ks.onIdleTimeout(ctx, connection)
	}
	return ctx
}

// OnClose 连接关闭后会执行的处理函数
func (ks *KServer) OnClose(connection netpoll.Connection) error {
	if !connection.IsActive() {
		if err := connection.Close(); err != nil {
			return err
		}
		log.Printf("[%s] 已经关闭连接.\n", connection.RemoteAddr().String())
	}
	return nil
}

// 连接空闲处理
// 当连接处于空闲一定时间后，主动将连接踢下线
func (ks *KServer) onIdleTimeout(ctx context.Context, connection netpoll.Connection) context.Context {
	// 连接空闲超时处理，由于目前netpoll不支持连接超时，所以需要手动实现
	// 创建一个channel，作为连接的心跳通道
	ch := make(chan struct{})
	// 开启一个goroutine，监听连接的心跳，一旦超过指定的时间未收到心跳，则关闭通道
	// TODO 后续使用协程池来维护
	go func() {
		for {
			select {
			case <-ch:
			case <-time.After(ks.idleTimeout):
				// 超时，关闭连接
				_ = connection.Close()
				close(ch)
			}
		}
	}()
	// 将心跳通道，在上下文中传递
	return context.WithValue(ctx,idleTimeoutChKey, ch)
}
