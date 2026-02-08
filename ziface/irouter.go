package ziface

// IRouter 配置当前连接的处理业务方法
type IRouter interface {
	PreHandle(request IRequest)  // 在处理conn业务之前的钩子方法
	Handle(request IRequest)     // 处理conn业务的方法
	PostHandle(request IRequest) // 处理conn业务之后的钩子方法
}
