// @Title normal_packer.go
// @Description
// @Author Zero - 2023/8/21 17:05:10

package packer

import (
	"encoding/binary"
	"github.com/zlx2019/kinx/kiface"
	"github.com/zlx2019/kinx/knet/message"
	"io"
)

const (
	// HeaderByteSize 消息内容长度所占字节数
	HeaderByteSize = 8
	// IDByteSize 消息ID所占字节数
	IDByteSize = 8

	// IDEndPos ID字段末尾字节位置
	IDEndPos = 16
)

// NormalPacker 消息数据包处理器: 根据固定的数据头长度进行解析,以 uint64(8byte)为准;
type NormalPacker struct {
	byteOrder binary.ByteOrder
}

// NewNormalPacker 构造函数
func NewNormalPacker() kiface.IPacker {
	return &NormalPacker{byteOrder: binary.BigEndian}
}

// Pack 消息打包
func (packer *NormalPacker) Pack(message kiface.IMessage) ([]byte, error) {
	// 计算数据包的总大(8 + 8 + 消息内容长度)
	totalSize := HeaderByteSize + IDByteSize + len(message.Payload())
	// 分配数据包缓冲区
	packs := make([]byte, totalSize)
	// 写入消息内容长度
	packer.byteOrder.PutUint64(packs[:HeaderByteSize], message.Len())
	// 写入消息ID
	packer.byteOrder.PutUint64(packs[HeaderByteSize:IDEndPos], message.ID())
	// 写入消息内容
	copy(packs[IDEndPos:], message.Payload())
	return packs, nil
}

// UnPack 消息解包
func (packer *NormalPacker) UnPack(reader io.Reader) (kiface.IMessage, error) {
	buf := make([]byte, HeaderByteSize+IDByteSize)
	// 读取消息内容长度和消息ID到缓冲区
	// 这里会阻塞读取，直到读取到指定长度的数据或者发生错误
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	// 解析内容长度和消息ID
	lens := packer.byteOrder.Uint64(buf[:HeaderByteSize])
	id := packer.byteOrder.Uint64(buf[HeaderByteSize:IDEndPos])
	// 读取消息内容
	payloadBuf := make([]byte, lens)
	_, err = io.ReadFull(reader, payloadBuf)
	if err != nil {
		return nil, err
	}
	return message.NewMessage(id, payloadBuf), nil
}
