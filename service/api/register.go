package api

import (
	"Rpc/decode"
	"Rpc/service/logic"
	"reflect"
	"time"
)

// 构造注册信息

func (s *Services) Register(v interface{}, weight ...int64) error {
	m := decode.CMS{
		OpenCode: 1,
		Message:  make(map[string]interface{}),
	}
	// 反射获取类型
	t := reflect.TypeOf(v)
	// 反射获取方法
	value := reflect.ValueOf(v)
	// 反射保存方法
	service := service{
		Method: make(map[string]reflect.Value),
	}
	// 对于不是指针的类型 无法反射
	if t.Kind() != reflect.Ptr {
		panic("err can't register a no ptr value ")
	}
	var methods []string
	t = t.Elem()
	// 获取方法
	for i := 0; i < t.NumMethod(); i++ {
		service.Method[t.Method(i).Name] = value.Method(i)
		methods = append(methods, t.Method(i).Name)
	}
	if len(weight) != 0 {
		m.Message["weight"] = weight[0]
	}
	s.Service = service
	// 对服务向远程发起注册 包含服务名 服务所包含的方法、
	//t.Name  是服务名
	m.Message[t.Name()] = methods
	// 请求 id 用于标识请求
	m.Id = time.Now().Unix() % 1000
	// 完成注册 进行监听
	err := logic.Process(s.con, m)

	return err
}
