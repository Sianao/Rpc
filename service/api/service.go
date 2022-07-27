package api

import (
	"net"
	"reflect"
)

// serivce  保存信息 一个节点之应该包含一个服务

type Services struct {
	con     net.Conn
	Service service
}

// 服务端信息 服务名 方法列表  可能存在多个
// 单个服务 的注册 信息
// 发起注册 可以将 服务的参数
type service struct {
	Method map[string]reflect.Value
}

func NewService(addr ...string) *Services {
	add := ":4563"
	if addr != nil {
		add = addr[0]
	}
	con, err := net.Dial("tcp", add)
	if err != nil {
		panic(err)
	}
	return &Services{
		con: con,
	}
}

// 将服务保存在本地 并向远程申请注册
