// @Title options.go
// @Description
// @Author Zero - 2023/8/30 19:16:46

package knet

import (
	"github.com/panjf2000/ants/v2"
	"github.com/zlx2019/kinx/kiface"
	"log"
	"time"
)

// OnHandler 注册服务的数据处理器
func (n *NormalServer) onHandler(handler kiface.IHandler) {
	n.handler = handler
}

// onOptions 注册服务的配置选项
func (n *NormalServer) onOptions(options ...NormalServerOption) {
	for _, option := range options {
		option(n)
	}
}

// setIdleTimeout 开启连接空闲超时,并且设置超时时间
func (n *NormalServer) setIdleTimeout(timeout time.Duration) {
	n.isIdleTimeout = true
	n.idleTimeout = timeout
}

// NormalServerOption NormalServer服务端的配置注册函数
type NormalServerOption func(server *NormalServer)

// WithHandler 设置处理器
func WithHandler(handler kiface.IHandler) NormalServerOption {
	return func(s *NormalServer) {
		s.onHandler(handler)
	}
}

// WithIdleTimeout 设置连接空闲超时时间
func WithIdleTimeout(timeout time.Duration) NormalServerOption {
	return func(s *NormalServer) {
		s.setIdleTimeout(timeout)
	}
}

// WithPool 初始化协程池，指定协程池的协程容量
func WithPool(capacity int) NormalServerOption {
	return func(s *NormalServer) {
		s.pool = newPool(capacity)
	}
}

// 默认的协程池配置
func newPool(capacity int) *ants.Pool {
	pool, _ := ants.NewPool(capacity, func(opt *ants.Options) {
		// 是否关闭回收空闲的work
		opt.DisablePurge = true
		// 回收空闲work的间隔。当DisablePurge为false时才生效
		// 如5 * time.Second 表示空闲5秒后的work会被回收掉
		opt.ExpiryDuration = time.Second * 3
		// 在初始化池时是否进行内存预分配。
		opt.PreAlloc = true
		//指定是否使用非阻塞模式执行任务。如果设置为true，则在协程池已满的情况下，任务会立即返回一个err，而不是等待空闲协程。
		// false表示不开启,阻塞等待可用的协程。
		opt.Nonblocking = false
		// 阻塞模式下,最多允许阻塞等待的协程数量。
		opt.MaxBlockingTasks = 100
		// 设置日志器
		opt.Logger = log.Default()
		// 指定一个函数用于处理协程中的 panic 异常。
		// TODO 暂时没有好的方案 不处理
		opt.PanicHandler = func(i interface{}) {
			log.Printf("ants pool panic: %v \n", i)
		}
	})
	pool.Running()
	return pool
}
