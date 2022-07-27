package api

import (
	"Rpc/decode"
	"Rpc/master/master"
	"errors"
	"net"

	"github.com/sirupsen/logrus"
)

var (
	ERR      = "err"
	NOMETHOD = "no such method"
)

// Call  函数调用 包含超时检测开启
func Call(c decode.CMS, con net.Conn) {
	// 解析请求中的参数 对其进行校验
	// 保存用户信息 方便下一次调用
	master.Job.Con[c.Id] = con
	service, ok := c.Message["Service"].(string)
	Method, ok := c.Message["Method"].(string)
	Param, ok := c.Message["Param"]
	if !ok {
		return
	}
	//加锁 保证并发安全
	master.Controller.Lock.Lock()
	defer master.Controller.Lock.Unlock()
	// 获取查看是否有该服务

	s, ok := master.Controller.Service[service]
	// 查看该服务的方法
	var target = false
	// 遍历 查看是否包含该方法
	if ok {
		for _, v := range s.Method {
			if v == Method {
				target = true
			}
		}
	}
	call := decode.CMS{
		Id:      c.Id,
		Message: make(map[string]interface{}),
	}
	//用于处理不存在的服务调用
	if !target || !ok {
		call.Message[ERR] = NOMETHOD
		CallBackErr(call, con)
		delete(master.Job.Con, c.Id)
		return
	}
	// 日志记录
	logrus.Info("call ", service, "'s method ", Method)
	call.OpenCode = 2
	call.Message["Method"] = Method
	call.Message["Param"] = Param
	code, err := decode.Encode(call)
	if err != nil {
		return
	}

	// 负载均衡选取最优服务器发送
	err = LoadBalance(s.Cons, code)
	if err != nil {
		CallBackErr(decode.CMS{
			OpenCode: 4,
			Id:       c.Id,
			Message: map[string]interface{}{
				"err": err.Error(),
			},
		}, con)
		delete(master.Job.Con, c.Id)
		return
	}

	// 启动超时检测 当服务端超过时间没返回数据的时候 响应错误
	go TimeOutCheck(c.Id)
	return
}

//Callback 函数返回
func Callback(c decode.CMS, con net.Conn) {
	bytes, _ := decode.Encode(c)
	_, _ = con.Write(bytes)
	con.Close()
}

func CallBackErr(msg decode.CMS, con net.Conn) {
	msg.OpenCode = 4
	//msg.Message["err"] = "错误 没有该服务或者方法"
	encode, _ := decode.Encode(msg)
	con.Write(encode)
	con.Close()
	return
}

// LoadBalance  负载均衡类似nginx 的权重计算 略有省略
func LoadBalance(cons []*master.Con, msg []byte) error {
	var (
		total int64 = 0
		min   int64 = 0
		ind         = -1
	)
	if len(cons) == 0 {
		return errors.New("no service can be use")
	}
	for k, v := range cons {
		if v.Active {
			v.CurrentWeight += v.Weight
			if v.CurrentWeight > min {
				ind = k
			}
			total += v.CurrentWeight
		}
	}
	if ind == -1 {
		return errors.New("no service can use")
	} else {
		write, err := cons[ind].Con.Write(msg)
		if err != nil {
			logrus.Error(err, write)
		}
		cons[ind].CurrentWeight -= total
		return err
	}

}
