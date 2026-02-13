package utils

import (
	"encoding/json"
	"github.com/Txinkang/zinx/ziface"
	"os"
)

type GlobalObj struct {
	/*
		server
	*/
	TcpServer ziface.IServer // 当前Zinx正在启动服务的Server对象
	Host      string         // 当前服务器主机IP
	TcpPort   int            // 当前服务器主机监听端口号
	Name      string         // 当前服务器名称

	/*
		zinx
	*/
	Version          string // 当前Zinx版本号
	MaxPacketSize    uint32 // 读取数据包的最大值
	MaxConn          int    // 当前服务器主机允许的最大连接个数
	WorkerPoolSize   uint32 // 工作池最大数量
	MaxWorkerTaskLen uint32 // 每条任务最大存储量

	/*
		config file path
	*/
	ConfFilePath string
}

// 定义一个全局对象
var GlobalObject *GlobalObj

// 读取用户配置文件
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// 解析json到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
}

// 提供init方法，默认加载
func init() {
	// 初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "1.0.0",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
		ConfFilePath:  "conf/zinx.json",

		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	// 从配置文件加载用户配置的参数
	GlobalObject.Reload()
}
