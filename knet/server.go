// @Title server.go
// @Description 服务端
// @Author Zero - 2023/8/10 19:53:21

package knet

import (
	"context"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
	"log"
	"net"
	"time"
)

// KServer 服务端,基于net_poll网络库实现的NIO模型服务端
type KServer struct {
	// 服务名称
	Name string
	// 服务端协议
	Protocol string
	// 服务端IP
	IP string
	// 服务端端口
	Port int

	// 服务端的网络侦听器
	listener net.Listener
	// 服务端的事件循环调度器，一个真正的NIO Server
	loop netpoll.EventLoop

	// 是否开启连接空闲超时机制，默认为false
	enableIdleTimeout bool
	// 连接空闲超时时间,默认为 30 min
	idleTimeout time.Duration
}

// NewServer 创建服务，返回一个IServer实例
// @param	name		服务名称
// @param	protocol	服务协议
func NewServer(name, protocol, ip string, port int) kiface.IServer {
	// 创建服务实例
	server := &KServer{
		Name:              name,
		IP:                ip,
		Protocol:          protocol,
		Port:              port,
		enableIdleTimeout: true,
		idleTimeout:       time.Second * 60 * 30,
	}
	return server
}

// Run 运行服务
func (ks *KServer) Run() error {
	// 准备服务
	if err := ks.ready(); err != nil {
		return err
	}
	// 开始运行服务，这里会阻塞，直到关闭服务
	if err := ks.start(); err != nil {
		return err
	}
	return nil
}

// 服务运行准备
func (ks *KServer) ready() error {
	// 建立TCP服务
	if err := ks.createNetwork(); err != nil {
		return err
	}
	// 创建事件循环调度器
	if err := ks.createEventLoop(); err != nil {
		return err
	}
	return nil
}

// 开始运行服务
func (ks *KServer) start() error {
	log.Printf("%s running successful... \n", ks.Name)
	// 阻塞运行调度器，直到调度器被关闭...
	if err := ks.loop.Serve(ks.listener); err != nil {
		return err
	}
	return nil
}

// 建立网络与事件循环组件
func (ks *KServer) createNetwork() error {
	address := fmt.Sprintf("%s:%d", ks.IP, ks.Port)
	// 建立一个TCP端点地址
	tcpAddr, err := net.ResolveTCPAddr(ks.Protocol, address)
	if err != nil {
		log.Println("resolve tcp addr failed cause: ", err.Error())
		return err
	}
	// 监听TCP服务端点
	ks.listener, err = net.ListenTCP(ks.Protocol, tcpAddr)
	if err != nil {
		log.Println("tcp listener failed cause: ", err.Error())
		return err
	}
	log.Println("tcp server ready successful. address in:", ks.listener.Addr().String())
	return nil
}

// 创建事件循环调度器，该调度器负责调度事件的处理函数
func (ks *KServer) createEventLoop() error {
	// 创建事件循环
	loop, err := netpoll.NewEventLoop(
		// 事件处理函数
		ks.OnRequest,
		// 连接初始化hook函数
		netpoll.WithOnPrepare(ks.OnPrepare),
		// 连接成功hook函数
		netpoll.WithOnConnect(ks.OnConnect),
		// 设置连接读取超时时间
		netpoll.WithReadTimeout(time.Second*3),
		// 设置连接写入超时时间
		netpoll.WithWriteTimeout(time.Second*3))
	if err != nil {
		log.Println("event loop create failed cause: ", err.Error())
		return err
	}
	log.Println("event loop ready successful.")
	ks.loop = loop
	return nil
}

// Shutdown 关闭服务器
func (ks *KServer) Shutdown() error {
	// 关闭事件循环
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return ks.loop.Shutdown(ctx)
}
