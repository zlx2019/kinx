// @Title server.go
// @Description 同步阻塞式服务端实现
// @Author Zero - 2023/8/21 13:35:03

package knet

import (
	"context"
	"flag"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/zlx2019/kinx/kiface"
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
	// 协程池
	pool *ants.Pool
	// 服务端的TCP服务监听器
	listener net.Listener
}

// NewNormalServer 创建服务端
// @param	opts	服务配置
func NewNormalServer(opts ...NormalServerOption) kiface.IServer {
	// 解析命令行参数
	flag.Parse()
	// 加载配置文件
	loadConfigs()
	server := &NormalServer{
		name:        configs.Name,
		protocol:    "tcp",
		iP:          configs.Host,
		port:        configs.Port,
		stopTrigger: make(chan struct{}),
	}
	// 注册要设置的配置
	server.onOptions(opts...)
	return server
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

	// 开启协程任务，开始接收客户端连接并且处理
	_ = n.pool.Submit(n.start)

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
	session := NewNormalSession(n.nextSessionID, conn, n.handler, ctx, nil, n.isIdleTimeout, n.idleTimeout)
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
		// 判断当前协程池内数量是否够用
		if !n.checkTaskQuantity() {
			_, _ = conn.Write([]byte("当前系统繁忙，请稍后再试~"))
			_ = conn.Close()
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
		session := NewNormalSession(n.nextSessionID, conn, n.handler, sessionCtx, cancel, n.isIdleTimeout, n.idleTimeout)

		// 启动3个协程，分别执行读、写任务以及心跳监控
		_ = n.pool.Submit(session.Reader)
		_ = n.pool.Submit(session.Writer)
		if n.isIdleTimeout {
			_ = n.pool.Submit(session.idleTimeOuter)
		}
		// ID自增
		n.nextSessionID++

		fmt.Printf("[%s] 会话运行成功，当前系统任务运行数量: %d \n", session.GetRemoteAddr(), n.pool.Running())
	}
}

// 查看当前可用的空闲协程是否足够
func (n *NormalServer) checkTaskQuantity() bool {
	if n.isIdleTimeout {
		return n.pool.Free() >= 3
	}
	return n.pool.Free() >= 2
}

// Shutdown 停止服务
func (n *NormalServer) Shutdown() error {
	if n.isRunning {
		// 关闭服务端
		n.stopTrigger <- struct{}{}
	}
	return nil
}
