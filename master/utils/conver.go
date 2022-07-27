package utils

import (
	"Rpc/decode"
	"github.com/sirupsen/logrus"
	"net"
)

// ParamService 解析服务
func ParamService(c map[string]interface{}) (service string, method []string, weight int64) {
	w, ok := c["weight"]
	if ok {
		weight, ok = w.(int64)
	} else {
		weight = 1
	}
	for k, v := range c {

		if k != "weight" {
			service = k
			// 这里保存了他所含有的方法 以及服务名
			value, _ := v.([]interface{})
			for _, v := range value {
				m, ok := v.(string)
				if ok {
					method = append(method, m)
				}
			}
			break
		}
	}

	return
}

// ReadFromCon 从连接中读取请求
func ReadFromCon(con net.Conn) (decode.CMS, error) {
	m := make([]byte, 1024)
	read, err := con.Read(m)
	if err != nil {
		logrus.Error(err)
		return decode.CMS{}, err
	}
	//解析 获得信息
	c := decode.Decode(m[0 : read+1])
	return c, nil
}

// WriteToCon 向连接中写信息
func WriteToCon(con net.Conn, msg decode.CMS) error {
	bytes, err := decode.Encode(msg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = con.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
