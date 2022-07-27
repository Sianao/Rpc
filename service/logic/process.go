package logic

import (
	"Rpc/decode"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
)

func Process(con net.Conn, m decode.CMS) error {
	encode, err := decode.Encode(m)
	if err != nil {
		return err
	}
	_, err = con.Write(encode)
	msg := make([]byte, 1024)
	n, err := con.Read(msg)
	//con.Close()
	c := decode.Decode(msg[0 : n+1])
	if value, ok := c.Message["err"]; ok {
		v, _ := value.(string)
		return errors.New(v)
	}
	logrus.Info("registering  service success")
	return err
}
