package logic

import (
	"Rpc/master/api"
	"Rpc/master/utils"
	"github.com/sirupsen/logrus"
	"net"
)

func Process(con net.Conn) {
	// 读取信息
	c, err := utils.ReadFromCon(con)
	if err != nil {
		// 处理错误
		logrus.Info("process err", err)
		con.Close()
		return
	}
	switch c.OpenCode {
	// 注册服务
	case 1:
		go api.Register(con, c)
		break
	//调用服务
	case 2:
		go api.Call(c, con)
		break
	}

}
