// @Title normal_server.go
// @Description 同步阻塞式服务端实现
// @Author Zero - 2023/8/21 13:35:03

package server

import (
	"context"
	"fmt"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/session"
	"net"
	"time"
)

// NormalServer 基础服务端,基于原生net库的同步阻塞的服务端
type NormalServer struct {
	// 服务名称
	name string
	// 服务端协议
	protocol string
	// 服务端IP
	iP string
	// 服务端端口
	port int
	// 下一个建立连接的会话ID，采用自增策略
	nextSessionID uint32
	// 服务是否处于启动状态
	isRunning bool
	// 会话是否开启空闲超时处理
	isIdleTimeout bool
	// 会话空闲超时时间，连接空闲超过该时间强制关闭
	idleTimeout time.Duration
	// 服务端关闭信号
	stopTrigger chan struct{}
	// 会话处理器
	handler kiface.IHandler
	// 服务端的TCP服务监听器
	listener net.Listener
}

// NewNormalServer 创建服务端
// @param	name	服务名
// @param	ip		服务IP
// @param	port	服务端口
// @param	opts	服务配置
func NewNormalServer(name, ip string, port int, opts ...kiface.BlockServerOption) kiface.IBlockServer {
	server := &NormalServer{
		name:        name,
		protocol:    "tcp",
		iP:          ip,
		port:        port,
		stopTrigger: make(chan struct{}),
	}
	// 注册要设置的配置
	server.OnOptions(opts...)
	return server
}

// OnHandler 注册服务的数据处理器
func (n *NormalServer) OnHandler(handler kiface.IHandler) {
	n.handler = handler
}

// OnOptions 注册服务的配置选项
func (n *NormalServer) OnOptions(options ...kiface.BlockServerOption) {
	for _, option := range options {
		option(n)
	}
}

// SetIdleTimeout 开启连接空闲超时,并且设置超时时间
func (n *NormalServer) SetIdleTimeout(timeout time.Duration) {
	n.isIdleTimeout = true
	n.idleTimeout = timeout
}

// Run 运行服务，并且阻塞监听连接
func (n *NormalServer) Run() error {
	// 创建TCP服务
	if err := n.ready(); err != nil {
		fmt.Println("tcp server ready failed cause: ", err.Error())
		return err
	}
	// 标记服务为运行状态
	n.isRunning = true
	fmt.Printf("%s running successful. address in: %s \n", n.name, n.listener.Addr().String())
	// 启动连接处理，开始接收客户端连接并且处理
	go n.start()

	//TODO 额外业务处理

	// 阻塞等待服务关闭
	select {
	case <-n.stopTrigger:
		fmt.Printf("%s shutodwn successful. \n", n.name)
	}
	return nil
}

// AsyncRun 异步运行服务，通过API来主动获取连接会话，进行处理
func (n *NormalServer) AsyncRun() error {
	// 创建TCP服务
	if err := n.ready(); err != nil {
		fmt.Println("tcp server ready failed cause: ", err.Error())
		return err
	}
	// 标记服务为运行状态
	n.isRunning = true
	fmt.Printf("%s running successful. address in: %s \n", n.name, n.listener.Addr().String())
	return nil
}

// GetSession 阻塞等待客户端连接，并且封装为会话
func (n *NormalServer) GetSession() (kiface.ISession, error) {
	conn, err := n.listener.Accept()
	if err != nil {
		return nil, err
	}
	var ctx context.Context
	if n.handler != nil {
		ctx = n.handler.OnConnectHandler(conn)
	}
	session := session.NewNormalSession(n.nextSessionID, conn, n.handler, ctx, nil, n.isIdleTimeout, n.idleTimeout)
	return session, nil
}

// 创建TCP网络服务
func (n *NormalServer) ready() error {
	if n.isRunning {
		panic("server already running")
	}
	address := fmt.Sprintf("%s:%d", n.iP, n.port)
	// 获取一个TCP的Addr
	tcpAddr, err := net.ResolveTCPAddr(n.protocol, address)
	if err != nil {
		return err
	}
	// 监听指定的Addr，获取监听器，
	n.listener, err = net.ListenTCP(n.protocol, tcpAddr)
	if err != nil {
		return err
	}
	return nil
}

// 异步循环处理客户端连接
func (n *NormalServer) start() {
	for {
		// 阻塞等待客户端连接
		conn, err := n.listener.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("Conn session successful. ID of: %d \n", n.nextSessionID)
		// 连接建立完成，回调连接建立事件处理函数，获取自定义的会话的上下文
		var ctx context.Context
		if n.handler != nil {
			ctx = n.handler.OnConnectHandler(conn)
		}
		// 创建会话的上下文，用于控制会话的退出
		sessionCtx, cancel := context.WithCancel(ctx)
		// 根据连接，创建一个连接会话，并且启动会话
		session.NewNormalSession(n.nextSessionID, conn, n.handler, sessionCtx, cancel, n.isIdleTimeout, n.idleTimeout).Rnu()
		// ID自增
		n.nextSessionID++
	}
}

// Shutdown 停止服务
func (n *NormalServer) Shutdown() error {
	if n.isRunning {
		// 关闭服务端
		n.stopTrigger <- struct{}{}
	}
	return nil
}
