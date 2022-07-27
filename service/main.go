package main

import (
	"Rpc/service/api"
	"github.com/sirupsen/logrus"
)

type MS struct{}

func (H MS) Add() string {
	return "hello"
}
func (MS) HAy() {

}
func main() {
	s := api.NewService(":4563")
	// 注册服务
	err := s.Register(&MS{}, 1)
	if err != nil {
		logrus.Error(err)
		return
	}
	s.Listen()
}
