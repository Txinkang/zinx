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
	//发送message数据到客户端
	SendMsg(msgId uint32, data []byte) error
}

// HandFunc 统一处理连接业务的接口,socket原声连接、请求数据、请求数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
