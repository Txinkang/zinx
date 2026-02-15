package ziface

import "net"

// IConnection 定义连接接口
type IConnection interface {
	// 启动连接，让当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态
	Stop()
	// 获取原始socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接ID
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	//发送message数据到客户端，无缓冲
	SendMsg(msgId uint32, data []byte) error
	//发送message数据到客户端，有缓冲
	SendBuffMsg(msgId uint32, data []byte) error
	// 设置连接属性
	SetProperty(key string, value interface{})
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
}

// HandFunc 统一处理连接业务的接口,socket原声连接、请求数据、请求数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
