package znet

import (
	"errors"
	"fmt"
	"github.com/Txinkang/zinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                  // 读写连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 对map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("conn add to ConnManager :conn num = ", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 对map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("conn remove, connId = ", conn.GetConnID(), ", conn num = ", connMgr.Len())
}

// 获取连接
func (connMgr *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 对map加读锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 获取连接
	if conn, ok := connMgr.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn not exist")
	}
}

// 获取当前连接数量
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 消除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 对map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 遍历所有连接挨个关闭
	for connId, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connId)
	}

	fmt.Println("conn clear conn num = ", connMgr.Len())

}
