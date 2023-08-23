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
	Name string
	// 服务端协议
	Protocol string
	// 服务端IP
	IP string
	// 服务端端口
	Port int

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
	// 读取事件处理函数
	onReadHandler kiface.OnHandler
	// 连接事件处理函数
	onConnectHandler kiface.OnConnectHandler
	// 服务端的TCP服务监听器
	listener net.Listener
}

// NewNormalServer 创建服务端
func NewNormalServer(name, protocol, ip string, port int, connectHandler kiface.OnConnectHandler, handler kiface.OnHandler) *NormalServer {
	return &NormalServer{
		Name:             name,
		Protocol:         protocol,
		IP:               ip,
		Port:             port,
		nextSessionID:    1,
		isRunning:        false,
		isIdleTimeout:    true,
		idleTimeout:      time.Second * 30,
		stopTrigger:      make(chan struct{}),
		onReadHandler:    handler,
		onConnectHandler: connectHandler,
	}
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
	fmt.Printf("%s running successful. address in: %s \n", n.Name, n.listener.Addr().String())
	// 启动连接处理，开始接收客户端连接并且处理
	go n.start()

	//TODO 额外业务处理

	// 阻塞等待服务关闭
	select {
	case <-n.stopTrigger:
		fmt.Printf("%s shutodwn successful. \n", n.Name)
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
	fmt.Printf("%s running successful. address in: %s \n", n.Name, n.listener.Addr().String())
	return nil
}

// GetSession 阻塞等待客户端连接，并且封装为会话
func (n *NormalServer) GetSession() (kiface.ISession, error) {
	conn, err := n.listener.Accept()
	if err != nil {
		return nil, err
	}
	var ctx context.Context
	if n.onConnectHandler != nil {
		ctx = n.onConnectHandler(conn)
	}
	session := session.NewNormalSession(n.nextSessionID, conn, n.onReadHandler, ctx, n.isIdleTimeout, n.idleTimeout)
	return session, nil
}

// 创建TCP网络服务
func (n *NormalServer) ready() error {
	if n.isRunning {
		panic("server already running")
	}
	address := fmt.Sprintf("%s:%d", n.IP, n.Port)
	// 获取一个TCP的Addr
	tcpAddr, err := net.ResolveTCPAddr(n.Protocol, address)
	if err != nil {
		return err
	}
	// 监听指定的Addr，获取监听器，
	n.listener, err = net.ListenTCP(n.Protocol, tcpAddr)
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
		// 回调连接事件处理函数，获取会话的上下文
		var ctx context.Context
		if n.onConnectHandler != nil {
			ctx = n.onConnectHandler(conn)
		}
		// 根据连接，创建一个连接会话，并且启动会话
		session.NewNormalSession(n.nextSessionID, conn, n.onReadHandler, ctx, n.isIdleTimeout, n.idleTimeout).Rnu()
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
