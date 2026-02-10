package znet

type Message struct {
	DataLen uint32
	Id      uint32
	Data    []byte
}

func NewMsgPackage(id uint32, data []byte) *Message {
	m := &Message{
		DataLen: uint32(len(data)),
		Id:      id,
		Data:    data,
	}
	return m
}

// 获取消息数据段长度
func (msg *Message) GetDataLen() uint32 { return msg.DataLen }

// 获取消息ID
func (msg *Message) GetMsgId() uint32 { return msg.Id }

// 获取消息内容
func (msg *Message) GetData() []byte { return msg.Data }

// 设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) { msg.DataLen = len }

// 设计消息ID
func (msg *Message) SetMsgId(msgId uint32) { msg.Id = msgId }

// 设计消息内容
func (msg *Message) SetData(data []byte) { msg.Data = data }
