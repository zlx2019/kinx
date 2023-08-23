// @Title server.go
// @Description
// @Author Zero - 2023/8/21 14:16:32

package main

import (
	"fmt"
	"github.com/zlx2019/kinx/knet/server"
	"time"
)

// NewNormalServer 案例二，基于传统模型开发

func main() {
	// 创建服务
	server := server.NewNormalServer("[kinx V1.0]", "tcp", "127.0.0.1", 9780, nil, nil)
	// 异步启动服务
	if err := server.AsyncRun(); err != nil {
		panic(err)
	}
	for {
		// 阻塞获取客户端连接
		session, _ := server.GetSession()
		go func() {
			// 读取客户端数据
			message, err := session.Read(time.Second * 3)
			if err != nil {
				fmt.Println(err)
			}

			// 写入客户端数据
			session.Write(message)
		}()
	}

}
