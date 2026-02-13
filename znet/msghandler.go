package znet

import (
	"fmt"
	"github.com/Txinkang/zinx/utils"
	"github.com/Txinkang/zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgId对应的Router处理方法
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// 启动一个Worker的工作流程
func (mh *MsgHandler) StartOneWorker(workId int, taskQueue chan ziface.IRequest) {
	fmt.Println("start workerId :", workId)

	// 不断等待队列中的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 开启工作池
func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 为当前任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前任务
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 将消息发送给消息队列
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 获取WorkerId
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	fmt.Println("connId = ", request.GetConnection().GetConnID(), ", msgId = ", request.GetMsgId(), ", send to workerId :", workerId)
	// 发送消息到对应的消息队列
	mh.TaskQueue[workerId] <- request
}
