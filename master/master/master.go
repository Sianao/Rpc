package master

import (
	"net"
	"sync"
)

// Center   保存服务段信息
type Center struct {
	// 所监听的地址
	Listen net.Listener
	Lock   sync.Mutex
	//map 保存注册信息 还在考虑是否加入负载均衡
	Service map[string]*Service
}

type Service struct {
	Cons   []*Con
	Method []string
}

/// 我在想 怎么写负载均衡

type Con struct {
	// 通过chanel 控制消息发送
	Con           net.Conn
	Weight        int64 //权重 默认为1
	Active        bool  //活跃数
	CurrentWeight int64
}

// Controller
var Controller *Center

// jobs 保存返回信息 用于在调用的时候进行处理

var Job *Jobs

type Jobs struct {
	Con map[int64]net.Conn
}
