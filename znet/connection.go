package znet

import (
	"fmt"
	"github.com/Txinkang/zinx/ziface"
	"net"
)

type Connection struct {
	// 当前连接的socket TCP 套接字
	Conn *net.TCPConn
	// 当前连接的ID，也可以称为SessionID，ID全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 该连接的处理方法
	handleAPI ziface.HandFunc
	// 告知该连接已经退出的channel
	ExitBuffChan chan bool
}

// NewConnection 创建连接方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callback_api,
		ExitBuffChan: make(chan bool, 1),
	}
}

// StartReader 处理conn读数据
func (c *Connection) StartReader() {
	fmt.Println("StartReader Goroutine is Running ...")
	defer fmt.Println(c.RemoteAddr().String(), "StartReader Goroutine is Stopped")
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error:", err)
			c.ExitBuffChan <- true
			continue
		}

		// 调用业务
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("Handle API Error:", err, "connId:", c.ConnID)
			c.ExitBuffChan <- true
			return
		}
	}
}

// 启动连接、让当前连接开始工作
func (c *Connection) Start() {
	// 开启业务
	go c.StartReader()

	// 一直阻塞，得到退出消息就退出
	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	// 如果当前连接已关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，则在此调用

	// 关闭Socket连接
	c.Conn.Close()

	// 通知该连接已关闭
	c.ExitBuffChan <- true

	// 关闭该连接所有管道
	close(c.ExitBuffChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
