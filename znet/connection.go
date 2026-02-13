package znet

import (
	"errors"
	"fmt"
	"github.com/Txinkang/zinx/ziface"
	"io"
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
	MsgHandler ziface.IMsgHandle
	// 告知该连接已经退出的channel
	ExitBuffChan chan bool
	// 读写分离消息通道
	msgChan chan []byte
}

// NewConnection 创建连接方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}
}

// StartReader 处理conn读数据
func (c *Connection) StartReader() {
	fmt.Println("StartReader Goroutine is Running ...")
	defer fmt.Println(c.RemoteAddr().String(), "StartReader Goroutine is Stopped")
	defer c.Stop()
	for {
		// 创建封包拆包的对象
		dp := NewDataPack()

		// 读取数据到buf
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read Error:", err)
			c.ExitBuffChan <- true
			continue
		}

		// 拆包，获取头信息
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("UnPack Error:", err)
			c.ExitBuffChan <- true
			continue
		}

		// 获取data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read Error:", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		// 获取当前客户端的Request
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 调用业务
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// StartWriter 处理conn写数据
func (c *Connection) StartWriter() {
	fmt.Println("StartWriter Goroutine is Running ...")
	defer fmt.Println(c.RemoteAddr().String(), "StartWriter Goroutine is Stopped")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Write Error:", err)
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}
func (c *Connection) Start() {
	// 开启读写业务
	go c.StartReader()
	go c.StartWriter()

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
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is Closed")
	}

	// 封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack Error:", err)
		return errors.New("pack Error")
	}

	// 发给读写channel
	c.msgChan <- msg

	return nil
}
