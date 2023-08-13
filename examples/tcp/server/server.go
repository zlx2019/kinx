// @Title server.go
// @Description  基于kinx框架，开发一个tcp服务端
// @Author Zero - 2023/8/12 15:46:09

package main

import (
	"github.com/zlx2019/kinx/knet"
)

func main() {
	iServer := knet.NewServer("[kinx V0.1]", "tcp", "127.0.0.1", 9898)
	go func() {
		if err := iServer.Run();err != nil {
			panic(err)
		}
	}()
	select {}
}