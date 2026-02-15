package ziface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动服务器方法
	Start()
	// Stop 停止服务器方法
	Stop()
	// Serve 开启业务服务方法
	Serve()
	// AddRouter 添加路由，让客户端自定义连接处理方法
	AddRouter(msgId uint32, router IRouter)
	// 获取连接管理
	GetConnMgr() IConnManager

	// 设置conn创建时的hook函数
	SetOnConnStart(func(connection IConnection))
	// 设置conn断开时的hook函数
	SetOnConnStop(func(connection IConnection))
	// 调用conn创建时的hook函数
	CallOnConnStart(conn IConnection)
	// 调用conn断开时的hook函数
	CallOnConnStop(conn IConnection)
}
