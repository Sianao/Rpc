package api

import (
	"Rpc/decode"
	"Rpc/master/master"
	"github.com/sirupsen/logrus"
	"time"
)

// Heartbeat 心跳检测 对服务的存活性进行检测
func Heartbeat(serv *master.Con) {
	t := time.NewTicker(time.Second * 15)
	// 防止内存泄露
	defer t.Stop()
	var ping decode.CMS
	ping.OpenCode = 3
	code, _ := decode.Encode(ping)
	for {
		select {
		case <-t.C:
			// 定时器 进行心跳检测
			_, err := serv.Con.Write(code)
			if err != nil {
				serv.Con.Close()
				// 服务活性修改
				logrus.Error("ping err : service doesn't work")
				serv.Active = false
				return
			}
			t.Reset(time.Second * 5)
		}
	}
}

// TimeOutCheck 超时检测
func TimeOutCheck(id int64) {
	t := time.NewTicker(time.Second * 20)
	select {
	case <-t.C:
		// 超时重传
		t.Stop()
		co, ok := master.Job.Con[id]
		if ok {
			CallBackErr(decode.CMS{
				OpenCode: 4,
				Id:       id,
				// 返回错误
				//我在想 要不要加超时重传
				Message: map[string]interface{}{
					"err": "time out ",
				}}, co)
			delete(master.Job.Con, id)
			return
		}
	}
}
