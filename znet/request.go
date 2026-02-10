package znet

import "github.com/Txinkang/zinx/ziface"

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好连接
	msg  ziface.IMessage    // 客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求消息的Id
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
