package logic

import (
	"Rpc/decode"
	"errors"
	"net"
)

func Call(c net.Conn, d decode.CMS) ([]interface{}, error) {
	defer c.Close()
	v, err := decode.Encode(d)
	_, err = c.Write(v)
	if err != nil {
		return nil, err
	}
	callback := make(chan []byte)
	go Wait(callback, c)
	var b []byte
	b = <-callback
	msg := decode.Decode(b)

	switch msg.OpenCode {
	// 对待正常返回
	case 2:
		re, ok := msg.Message["Return"]
		if ok {
			ret, _ := re.([]interface{})
			return ret, nil
		}

	case 4:
		re, ok := msg.Message["err"]
		if ok {
			re, _ := re.(string)
			return nil, errors.New(re)
		}
		break
	}
	return nil, nil
}
