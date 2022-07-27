package cmd

import (
	"Rpc/master/logic"
	"Rpc/master/master"
	"github.com/sirupsen/logrus"
	"net"
)

func MasterInit(port *string) {
	master.Controller = &master.Center{}
	master.Job = &master.Jobs{}
	lister, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		panic("can't listen on port " + *port)
	}
	master.Controller.Listen = lister
	master.Job.Con = make(map[int64]net.Conn)
	master.Controller.Service = make(map[string]*master.Service)

}

// Start 开始注册服务发现
func Start(port *string) {
	// 初始化
	MasterInit(port)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
	})
	logrus.Info("Master is listening on port: ", *port)
	for {
		// 接受连接
		con, err := master.Controller.Listen.Accept()
		if err != nil {
			return
		}
		logic.Process(con)
	}
}
