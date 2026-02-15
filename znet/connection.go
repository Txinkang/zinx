package znet

import (
	"errors"
	"fmt"
	"github.com/Txinkang/zinx/utils"
	"github.com/Txinkang/zinx/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// 当前连接属于哪个Server
	TcpServer ziface.IServer
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
	// 读写分离消息通道，无缓冲管道
	msgChan chan []byte
	// 读写分离消息通道，有缓冲管道
	msgBuffChan chan []byte
	// 连接属性
	property map[string]interface{}
	// 设置属性锁
	propertyLock sync.RWMutex
}

// NewConnection 创建连接方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	// 初始化conn
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxWorkerTaskLen),
		property:     make(map[string]interface{}),
	}
	// 把当前conn加入到连接管理中
	c.TcpServer.GetConnMgr().Add(c)

	return c
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
			return
		}

		// 拆包，获取头信息
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("UnPack Error:", err)
			c.ExitBuffChan <- true
			return
		}

		// 获取data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read Error:", err)
				c.ExitBuffChan <- true
				return
			}
		}
		msg.SetData(data)

		// 获取当前客户端的Request
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 调用业务，开启工作池了就走工作池逻辑，没开启就还是临时goroutine逻辑
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 处理conn写数据
func (c *Connection) StartWriter() {
	fmt.Println("StartWriter Goroutine is Running ...")
	defer fmt.Println(c.RemoteAddr().String(), "StartWriter Goroutine is Stopped")

	for {
		select {
		// 无缓冲
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data Error:", err)
				return
			}
		// 有缓冲
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send buff data Error:", err)
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
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

	// 创建连接时执行的钩子方法
	c.TcpServer.CallOnConnStart(c)
}
func (c *Connection) Stop() {
	// 如果当前连接已关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 断开连接时执行的钩子方法
	c.TcpServer.CallOnConnStop(c)

	// 关闭Socket连接
	c.Conn.Close()

	// 通知该连接已关闭
	c.ExitBuffChan <- true

	// 从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭该连接所有管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
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
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is Closed")
	}

	// 封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack Error msg id = :", msgId)
		return errors.New("pack Error")
	}

	// 发给读写channel
	c.msgBuffChan <- msg

	return nil
}
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
