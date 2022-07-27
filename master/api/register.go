package api

import (
	"Rpc/decode"
	"Rpc/master/master"
	"Rpc/master/utils"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
)

const (
	DifferentSMethod = "registered service with different method"
	errString        = "err"
)

// 当注册完成后 则需要对信息进行轮询

func Register(conn net.Conn, c decode.CMS) {
	// 对服务以及方法 权重 进行提取
	service, method, weight := utils.ParamService(c.Message)
	// 构造返回消息
	res := decode.CMS{
		OpenCode: 1,
		Id:       c.Id,
		Message:  make(map[string]interface{})}
	// 加锁 保证并发安全
	//master.Controller.Lock.Lock()
	//defer master.Controller.Lock.Unlock()
	v, ok := master.Controller.Service[service]
	// 对于已经查到的服务 查看发送的注册请求是否符合规则 既 方法是否相同 进行检测
	con := master.Con{
		Con:           conn,
		Active:        true,
		Weight:        weight,
		CurrentWeight: -1,
	}
	if ok {
		// 检查方法是否一致
		if err := CheckMethod(v, method); err != nil {
			res.Message[errString] = err.Error()
			utils.WriteToCon(conn, res)
			conn.Close()
			return
		}
		// 查看是否有挂掉的服务 进行重新加载
		var tag = false
		for l, co := range v.Cons {
			if co.Active == false {
				v.Cons[l] = &con
				tag = true
				break
			}
		}
		// 服务恢复以及服务扩容
		if !tag {
			v.Cons = append(v.Cons, &con)
			logrus.Info("a service append ", service)
		} else {
			logrus.Info("a service recover ", service)
		}

		utils.WriteToCon(conn, res)
		go ServiceListen(&con)
		//go ServiceAccept(&c)
		go Heartbeat(&con)
	} else {
		// 添加服务 启动十个 进行
		c := master.Con{
			Con:    conn,
			Active: true,
			Weight: weight,
		}
		serv := master.Service{
			Cons:   []*master.Con{&c},
			Method: method,
		}
		logrus.Info("a service registered ", service)
		utils.WriteToCon(conn, res)
		master.Controller.Service[service] = &serv
		go ServiceListen(&c)
		//go ServiceAccept(&c)
		go Heartbeat(&c)
	}
	return
}

// CheckMethod 检查是否方法参数一致
func CheckMethod(v *master.Service, method []string) error {
	for k, m := range v.Method {
		if method[k] != m {
			// 说明注册的方法存在问题 不是同一个服务 却存在相同的服务名
			// 关闭连接 释放资源
			return errors.New(DifferentSMethod)
		}
	}
	return nil

}
