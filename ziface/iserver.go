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
}
