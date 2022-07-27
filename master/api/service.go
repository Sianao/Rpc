package api

import (
	"Rpc/master/master"
	"Rpc/master/utils"
	"github.com/sirupsen/logrus"
)

// ServiceListen 监听服务端发来的消息 对于一个已经注册的服务 只有回应消息 同时对服务进行检测
func ServiceListen(c *master.Con) {
	for {
		cms, err := utils.ReadFromCon(c.Con)
		if err != nil {
			c.Active = false
			logrus.Error("a service can't work")
			return
		}
		conn, ok := master.Job.Con[cms.Id]
		if ok {
			// 回调 返回信息
			Callback(cms, conn)
			delete(master.Job.Con, cms.Id)
		}
	}
}
