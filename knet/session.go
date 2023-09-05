// @Title session.go
// @Description
// @Author Zero - 2023/8/21 14:24:27

package knet

import (
	"context"
	"fmt"
	"github.com/zlx2019/kinx/kiface"
	"io"
	"net"
	"time"
)

// NormalSession 同步阻塞式客户端会话连接，用于管理客户端的连接，搭配NormalServer服务端使用;
type NormalSession struct {
	// 会话ID
	ID uint32
	// 客户端连接
	Conn net.Conn
	// 会话连接是否关闭
	IsClosed bool
	// 会话上下文
	context context.Context
	// 会话上下文取消方法
	cancel context.CancelFunc

	// 会话是否开启空闲超时处理
	isIdleTimeout bool
	// 会话空闲超时时间，连接空闲超过该时间强制关闭
	idleTimeout time.Duration

	// 会话处理器
	handler kiface.IHandler
	// 连接心跳通道，读取到连接的数据后刷新一下心跳，表示处于活跃
	active chan struct{}
	// 消息输出通道，将要发送给本会话的数据添加到该通道内，由写协程读取并且发送给连接
	outChannel chan kiface.IMessage
	// 消息封包与解包处理器
	packer kiface.IPacker
}

// NewNormalSession 创建连接会话
func NewNormalSession(id uint32, conn net.Conn, handler kiface.IHandler, ctx context.Context, cancel context.CancelFunc, isIdleTimeout bool, idleTimeout time.Duration) *NormalSession {
	return &NormalSession{
		ID:            id,
		Conn:          conn,
		IsClosed:      false,
		handler:       handler,
		isIdleTimeout: isIdleTimeout,
		idleTimeout:   idleTimeout,
		context:       ctx,
		cancel:        cancel,
		outChannel:    make(chan kiface.IMessage, 16),
		packer:        NewNormalPacker(),
	}
}

// Rnu 启动会话
func (ns *NormalSession) Rnu() {
	// 启动3个协程，分别执行读、写任务以及心跳监控
	go ns.Reader()
	go ns.Writer()
	if ns.isIdleTimeout {
		ns.active = make(chan struct{})
		go ns.idleTimeOuter()
	}
}

// Reader 连接会话的读任务,读取连接的数据，回调 onHandler 函数进行处理
func (ns *NormalSession) Reader() {
	fmt.Printf("[%s] Session ID: %d Reader Work Running... \n", ns.GetRemoteAddr(), ns.ID)
	// 循环读取数据
	for {
		// 阻塞读取消息数据，直到:读取到足够的数据 | 读取超时 | 连接被关闭
		message, err := ns.Read(time.Second * 3)
		// 读取错误处理
		if err != nil {
			if err == io.EOF || ns.IsClosed {
				// err == io.EOF 表示客户端主动关闭;
				// ns.IsClosed  表示服务端主动关闭，连接被close: 超时被强制关闭 | 处理函数抛出错误;
				fmt.Printf("[%s] Session ID: %d Reader Work Shutdown... \n", ns.GetRemoteAddr(), ns.ID)
				// 停止任务
				ns.Stop()
				return
			} else if e, ok := err.(net.Error); ok && e.Timeout() {
				// 本次读取数据超时
				continue
			}
			// TODO 其他错误处理
			continue
		}
		// 读取到会话连接的数据，回调注册的处理函数链
		if ns.handler != nil {
			ctx := NewHandlerContext(ns, message, ns.context)
			if err := ns.handler.OnHandler(ctx); err != nil {
				ns.Stop()
			}
		}
	}
}

// Writer 连接会话的写任务,读取会话的 outChannel 通道数据，将其写到客户端连接中.
func (ns *NormalSession) Writer() {
	fmt.Printf("[%s] Session ID: %d Writer Work Running... \n", ns.GetRemoteAddr(), ns.ID)
	for {
		select {
		case msg, ok := <-ns.outChannel:
			if !ok {
				// 消息通道已关闭，表示已经执行了Stop()方法，会话关闭
				fmt.Printf("[%s] Session ID: %d Writer Work Shutdown... \n", ns.GetRemoteAddr(), ns.ID)
				return
			}
			_ = ns.Write(msg)
		}
	}
}

// Send 将消息添加至会话通道，然后被写入到客户端连接中
func (ns *NormalSession) Send(message kiface.IMessage) {
	ns.outChannel <- message
}

// idleTimeOuter 会话的心跳检测器，超过指定时间未接到心跳则超时
func (ns *NormalSession) idleTimeOuter() {
	fmt.Printf("[%s] Session ID: %d Timeouter Running... \n", ns.GetRemoteAddr(), ns.ID)
	defer fmt.Printf("[%s] Session ID: %d Timeouter Shutdown... \n", ns.GetRemoteAddr(), ns.ID)
	for {
		select {
		case <-ns.active:
			// 接收心跳信号，保持活跃
		case <-ns.context.Done():
			// 客户端已关闭，停止心跳发送
			return
		case <-time.After(ns.idleTimeout):
			// 会话连接超时退出
			_, _ = ns.GetConn().Write([]byte("您超时了!"))
			ns.Stop()
			return
		}
	}
}

// 从会话连接中读取数据，并且解包
func (ns *NormalSession) Read(timeout time.Duration) (kiface.IMessage, error) {
	// 设置本次读取数据的阻塞超时时间 3s
	_ = ns.Conn.SetReadDeadline(time.Now().Add(timeout))
	// 从连接中阻塞读取数据，并且解包为IMessage
	return ns.packer.UnPack(ns.Conn)
}

// Write 向客户端连接写入数据
func (ns *NormalSession) Write(message kiface.IMessage) error {
	// 将消息体进行封包
	pack, err := ns.packer.Pack(message)
	if err != nil {
		return err
	}
	// 写入连接
	_, err = ns.Conn.Write(pack)
	return err
}

// GetConn 获取会话的客户端连接
func (ns *NormalSession) GetConn() net.Conn {
	return ns.Conn
}

// GetSessionID 获取会话的ID
func (ns *NormalSession) GetSessionID() uint32 {
	return ns.GetSessionID()
}

// GetRemoteAddr 获取客户端连接地址
func (ns *NormalSession) GetRemoteAddr() net.Addr {
	return ns.Conn.RemoteAddr()
}

// Stop 关闭会话
func (ns *NormalSession) Stop() {
	// 防止重复关闭两次，因为如果是连接超时关闭，会执行两次该函数.
	if !ns.IsClosed {
		// 将会话标记为已关闭
		ns.IsClosed = true
		// 关闭会话的消息通道，借此关闭写协程
		close(ns.outChannel)
		// 关闭会话上下文
		ns.cancel()
		// 执行 连接关闭的回调函数
		_ = ns.handler.OnClosedHandler(ns.Conn)
		// 关闭连接
		_ = ns.Conn.Close()
	}
}

// IsClose 会话是否已关闭
func (ns *NormalSession) IsClose() bool {
	return ns.IsClosed
}

// GetContext 获取会话的上下文
func (ns *NormalSession) GetContext() context.Context {
	return ns.context
}
