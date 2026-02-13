package znet

import (
	"fmt"
	"github.com/Txinkang/zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个MsgId对应的Router处理方法
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 以非阻塞的方式处理消息
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 获取对应消息的router
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId :", request.GetMsgId(), " not exist")
		return
	}

	// 执行
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为对应的消息添加具体的路由
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("api msgId is exist = " + strconv.Itoa(int(msgId)))
		return
	}
	mh.Apis[msgId] = router
	fmt.Println("add api msgId :", msgId)
}
