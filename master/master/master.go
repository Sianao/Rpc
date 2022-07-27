package master

import (
	"net"
	"sync"
)

// Controller 控制中心
var Controller *Center

// Job 保存返回信息 用于在调用的时候进行处理
var Job *Jobs

// Center  保存中心信息
type Center struct {
	Listen net.Listener
	Lock   sync.Mutex
	//map 保存注册信息
	Service map[string]*Service
}

// Service 单个服务信息
type Service struct {
	Cons   []*Con
	Method []string
}

// Con 单一连接的控制
type Con struct {
	Con           net.Conn
	Weight        int64 //权重 默认为1
	Active        bool  //活跃
	CurrentWeight int64 // 负载均衡
}

type Jobs struct {
	Con map[int64]net.Conn
}
