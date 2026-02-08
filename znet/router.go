package znet

import "github.com/Txinkang/zinx/ziface"

// BaseRouter 类为一切实现Router子类的父类。继承后，有不希望实现的方法也可以实例化。
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (br *BaseRouter) Handle(req ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
