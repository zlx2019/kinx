// @Title handler.go
// @Description  服务端事件处理
// @Author Zero - 2023/8/12 20:20:58

package knet

import (
	"context"
	"fmt"
	"github.com/cloudwego/netpoll"
	"log"
	"strings"
)

// OnRequest 事件处理函数
// 当一个连接事件就绪后，EventLoop会回调该函数
// @param	ctx		连接上下文
// @param	conn	事件就绪的连接
func (ks *KServer) OnRequest(ctx context.Context, conn netpoll.Connection) error {
	// 连接是否有效
	if !conn.IsActive() {
		log.Println("the connection is close.")
		return nil
	}
	// 获取连接读取器
	reader := conn.Reader()
	defer reader.Release()
	// 将读取器中的数据全部读出
	bytes, err := reader.Next(reader.Len())
	if err != nil {
		if err == netpoll.ErrReadTimeout {
			log.Printf("[%s] 读取数据超时 \n",conn.RemoteAddr().String())
		}
		return err
	}
	// 截取末尾的\n
	message := strings.TrimSuffix(string(bytes), "\n")
	fmt.Println(message)

	// 获取心跳通道，发送心跳，维持连接的活跃
	if idelCh,ok := ctx.Value(idleTimeoutChKey).(chan struct{}); ok{
		idelCh <- struct{}{}
	}
	return nil
}
