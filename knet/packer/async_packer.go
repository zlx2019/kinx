// @Title async_packer.go
// @Description	netPoll数据包处理器
// @Author Zero - 2023/8/22 17:18:37

package packer

import (
	"encoding/binary"
	"errors"
	"github.com/cloudwego/netpoll"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/message"
	"io"
	"unsafe"
)

var (
	// ErrWait 该错误表示连接内数据包不完整，需要等待更多数据
	ErrWait = errors.New("message payload incomplete")
)

// AsyncPacker AsyncServer非阻塞模型服务封包与拆包处理器
type AsyncPacker struct {
	byteOrder binary.ByteOrder
}

// NewAsyncPacker 创建非阻塞封包与拆包处理器
func NewAsyncPacker() kiface.IPacker {
	return &AsyncPacker{byteOrder: binary.BigEndian}
}

// Pack 消息打包
func (packer *AsyncPacker) Pack(message kiface.IMessage) ([]byte, error) {
	// 获取消息内容长度、消息ID、消息内容
	length, id, payload := message.Len(), message.ID(), message.Payload()
	// 计算缓冲区所需大小，并且分配
	totalSize := int(unsafe.Sizeof(length)) + int(unsafe.Sizeof(id)) + len(payload)
	packs := make([]byte, totalSize)
	// 写入内容长度
	packer.byteOrder.PutUint64(packs[:8], length)
	// 写入消息ID
	packer.byteOrder.PutUint64(packs[8:16], id)
	// 写入数据内容
	copy(packs[16:], payload)
	return packs, nil
}

// UnPack 数据解包
func (packer *AsyncPacker) UnPack(reader io.Reader) (kiface.IMessage, error) {
	// 强转为netPoll的连接
	if conn, ok := reader.(netpoll.Connection); ok {
		// 获取连接的异步读取器
		reader := conn.Reader()
		// 判断可读数据是否足够解析出 消息内容长度
		if reader.Len() < 8 {
			// 可读取数据不足，无法读取消息的内容长度
			return nil, ErrWait
		}
		// 读取消息内容长度,只是读取，并不会将数据弹出
		lengthBytes, _ := reader.Peek(8)
		length := packer.byteOrder.Uint64(lengthBytes)
		// 判断可读数据是否足够一个数据包
		if uint64(reader.Len()) < 8+8+length {
			return nil, ErrWait
		}
		// 跳过已读取的内容长度
		_ = reader.Skip(8)
		// 获取消息ID
		idBytes, _ := reader.Next(8)
		id := packer.byteOrder.Uint64(idBytes)
		// 根据读取到的内容长度数值，读取消息内容
		payload, _ := reader.Next(int(length))
		// 释放读取器缓冲区
		_ = reader.Release()
		return message.NewMessage(id, payload), nil
	}
	return nil, errors.New("reader type is not netPoll Connection")
}
