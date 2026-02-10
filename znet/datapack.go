package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Txinkang/zinx/utils"
	"github.com/Txinkang/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	// Id uint32 + Data uint32 共8字节
	return 8
}

// 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放Bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
func (dp *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	// 创建一个输入二进制数据的ioReader
	dataBuff := bytes.NewReader(data)

	// 只解压head的信息，获取dataLen和msgId
	msg := &Message{}

	// 获取dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 获取msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断长度是否合法
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("data too large")
	}

	return msg, nil
}
