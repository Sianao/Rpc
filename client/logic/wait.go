package logic

import (
	"net"
)

func Wait(callback chan []byte, c net.Conn) {

	for {
		bs := make([]byte, 1024)
		read, err := c.Read(bs)
		if err != nil {
			return
		}
		callback <- bs[0 : read+1]
	}
}
