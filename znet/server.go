package znet

import (
	"fmt"
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
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[START] server listenner at IP: %s, Port: %d, is starting\n", s.IP, s.Port)

	// 开启一个go去做服务器端Linster业务
	go func() {
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

		// 3、启动server网络连接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept err:", err)
				continue
			}
			// TODO 3.2 设置服务器最大连接控制，如果超过最大连接，则关闭新的连接

			// TODO 3.3 处理该新连接请求的业务方法，此时 handler 和 conn 应该是绑定的

			// 暂时做一个最大512字节的回显服务
			go func() {
				// 不断的循环，从客户端获取数据
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("read buf err:", err)
						continue
					}
					// 回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err:", err)
						continue
					}
				}
			}()
		}
	}()
}

// 关闭网络服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server stop, name", s.Name)

	// TODO 将需要清理的连接信息或其他信息一并停止或清理
}

// 开启业务服务
func (s *Server) Serve() {
	s.Start()

	// TODO 如果在启动服务的时候还要处理其他的事情，则可以在这里添加

	// 阻塞，否则主Go退出，listenner的go将会退出
	select {}
}

// 创建服务器句柄
func NewServe(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
