package znet

import (
	"fmt"
	"github.com/Txinkang/zinx/utils"
	"github.com/Txinkang/zinx/ziface"
	"net"
)

// IServer接口实现，定义一个Server服务类
type Server struct {
	// 服务器的名称
	Name string
	// tcp4 or other
	IPVersion string
	// 服务绑定的IP地址
	IP string
	// 服务绑定的端口
	Port int
	// 消息管理模块，用来绑定msgId和对应的处理方法
	msgHandler ziface.IMsgHandle
	// 连接管理
	ConnMgr ziface.IConnManager

	// 当前server创建连接时
	OnConnStart func(conn ziface.IConnection)
	// 当前server创建断开时
	OnConnStop func(conn ziface.IConnection)
}

// 创建服务器句柄
func NewServe() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		msgHandler:  NewMsgHandler(),
		ConnMgr:     NewConnManager(),
		OnConnStart: func(conn ziface.IConnection) {},
		OnConnStop:  func(conn ziface.IConnection) {},
	}
	return s
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[START]Server name:%s,listenner at IP: %s, Ports %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx]Version: %s,MaxConn: %d, MaxPacketSize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)

	// 开启一个go去做服务器端Linster业务
	go func() {
		// 0、开启工作池
		s.msgHandler.StartWorkerPool()
		// 1、获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		// 2、监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err:", err)
			return
		}
		// 已经成功监听
		fmt.Println("start zinx server", s.Name, "success, now listenning")

		// TODO 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		// 3、启动server网络连接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept err:", err)
				continue
			}
			// TODO 3.2 设置服务器最大连接控制，如果超过最大连接，则关闭新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			// TODO 3.3 处理该新连接请求的业务方法
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			// 3.4 启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

// 关闭网络服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server stop, name", s.Name)

	// 将需要清理的连接信息或其他信息一并停止或清理
	s.ConnMgr.ClearConn()
}

// 开启业务服务
func (s *Server) Serve() {
	s.Start()

	// TODO 如果在启动服务的时候还要处理其他的事情，则可以在这里添加

	// 阻塞，否则主Go退出，listenner的go将会退出
	select {}
}

// 注册路由业务方法，供客户端连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)

	fmt.Println("[START] AddRouter success")
}

// 获取连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 设置创建连接时的hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置断开连接时的hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用创建连接时的hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("[CallOnConnStart] s.OnConnStart is nil", s.Name)
		s.OnConnStart(conn)
	}
}

// 调用断开连接时的hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("[CallOnConnStop] s.OnConnStop is nil", s.Name)
		s.OnConnStop(conn)
	}
}
