// @Title async_server.go
// @Description 同步非阻塞服务端实现
// @Author Zero - 2023/8/10 19:53:21

package server

import (
	"context"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
	"log"
	"net"
	"time"
)

// AsyncServer 服务端 基于net_poll网络库实现的NIO模型服务端
type AsyncServer struct {
	// 服务名称
	Name string
	// 服务端协议
	Protocol string
	// 服务端IP
	IP string
	// 服务端端口
	Port int

	onPrepare kiface.OnAsyncPrepareHandler // 连接初始化事件
	onConnect kiface.OnAsyncConnectHandler // 连接完成事件
	onRead    kiface.OnAsyncReadHandler    // 连接读取事件
	onClosed  kiface.OnAsyncClosedHandler  // 连接关闭事件

	nextSessionID uint32
	// 连接是否开启空闲超时处理，空闲超时则强制关闭连接，默认为false;
	isIdleTimeout bool
	// 连接空闲超时时间,默认为 30 min
	idleTimeout time.Duration
	// 服务端的网络侦听器
	listener net.Listener
	// 服务端的事件循环调度器，一个真正的NIO Server
	loop netpoll.EventLoop
}

// OnOptions 注册服务端配置
func (as *AsyncServer) OnOptions(opts ...kiface.AsyncServerOption) {
	for _, option := range opts {
		option(as)
	}
}

// OnPrepare 注册[连接初始化事件]处理函数
func (as *AsyncServer) OnPrepare(prepareHandler kiface.OnAsyncPrepareHandler) {
	as.onPrepare = prepareHandler
}

// OnConnect 注册[连接完成事件]处理函数
func (as *AsyncServer) OnConnect(connectHandler kiface.OnAsyncConnectHandler) {
	as.onConnect = connectHandler
}

// OnReadHandler 注册[连接读取事件]处理函数
// 当连接中有可读数据时，会回调 onRead 函数
func (as *AsyncServer) OnReadHandler(readHandler kiface.OnAsyncReadHandler) {
	as.onRead = readHandler
}

// OnClosed 注册[连接关闭事件]处理函数
func (as *AsyncServer) OnClosed(closedHandler kiface.OnAsyncClosedHandler) {
	as.onClosed = closedHandler
}

// NewAsyncServer 创建NewNIOServer服务
// @param	name		服务名称
// @param	protocol	服务协议
func NewAsyncServer(name, protocol, ip string, port int) kiface.IAsyncServer {
	// 创建服务实例
	s := &AsyncServer{
		Name:          name,
		IP:            ip,
		Protocol:      protocol,
		Port:          port,
		nextSessionID: 0,
		isIdleTimeout: true,
		idleTimeout:   time.Second * 10,
	}
	return s
}

// Run 运行服务
func (as *AsyncServer) Run() error {
	// 准备服务
	if err := as.ready(); err != nil {
		return err
	}
	// 开始运行服务，这里会阻塞，直到关闭服务
	if err := as.start(); err != nil {
		return err
	}
	return nil
}

// 服务运行准备
func (as *AsyncServer) ready() error {
	// 建立TCP服务
	if err := as.createNetwork(); err != nil {
		return err
	}
	// 创建事件循环调度器
	if err := as.createEventLoop(); err != nil {
		return err
	}
	return nil
}

// 开始运行服务
func (as *AsyncServer) start() error {
	log.Printf("%s running successful... \n", as.Name)
	// 阻塞运行调度器，直到调度器被关闭...
	if err := as.loop.Serve(as.listener); err != nil {
		return err
	}
	log.Printf("%s shutdown successful... \n", as.Name)
	return nil
}

// 建立网络与事件循环组件
func (as *AsyncServer) createNetwork() error {
	address := fmt.Sprintf("%s:%d", as.IP, as.Port)
	// 建立一个TCP端点地址
	tcpAddr, err := net.ResolveTCPAddr(as.Protocol, address)
	if err != nil {
		log.Println("resolve tcp addr failed cause: ", err.Error())
		return err
	}
	// 监听TCP服务端点
	as.listener, err = net.ListenTCP(as.Protocol, tcpAddr)
	if err != nil {
		log.Println("tcp listener failed cause: ", err.Error())
		return err
	}
	log.Println("tcp server ready successful. address in:", as.listener.Addr().String())
	return nil
}

// 创建事件循环调度器，该调度器负责调度事件的处理函数
func (as *AsyncServer) createEventLoop() error {
	// 创建事件循环

	loop, err := netpoll.NewEventLoop(
		// 连接有数据可读事件处理
		as.onReadEvent,
		// 连接初始化hook函数
		netpoll.WithOnPrepare(as.onPrepareEvent),
		// 连接成功hook函数
		netpoll.WithOnConnect(as.onConnectEvent),
		// 设置连接读取超时时间
		netpoll.WithReadTimeout(time.Second*3),
		// 设置连接写入超时时间
		netpoll.WithWriteTimeout(time.Second*3))
	if err != nil {
		log.Println("event loop create failed cause: ", err.Error())
		return err
	}
	log.Println("event loop ready successful.")
	as.loop = loop
	return nil
}

// Shutdown 关闭服务器
func (as *AsyncServer) Shutdown() error {
	// 关闭事件循环
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return as.loop.Shutdown(ctx)
}
