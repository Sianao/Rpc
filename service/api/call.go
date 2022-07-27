package api

import (
	"Rpc/decode"
	"Rpc/master/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"reflect"
)

func (s *Services) Call(c decode.CMS, con net.Conn) {
	defer func(cms decode.CMS) {
		if l := recover(); l != nil {
			logrus.Error("call the service err", l)
		}
		utils.WriteToCon(s.con, decode.CMS{
			OpenCode: 4, Id: c.Id,
			Message: map[string]interface{}{
				"err": "request bad no support param or method",
			}})
	}(c)

	v, _ := c.Message["Method"]
	param, _ := c.Message["Param"]
	prams, ok := param.([]interface{})
	var m []reflect.Value
	if ok {
		for _, v := range prams {
			m = append(m, reflect.ValueOf(v))
		}
	}
	method, _ := v.(string)
	fmt.Println("call the method", method, " with param ", param)
	result := s.Service.Method[method].Call(m)
	var re []interface{}
	for _, v := range result {
		re = append(re, v.Interface())
	}
	r := decode.CMS{
		OpenCode: 2,
		Id:       c.Id,
		Message:  make(map[string]interface{}),
	}
	r.Message["Return"] = re

	b, _ := decode.Encode(r)
	s.con.Write(b)
}
