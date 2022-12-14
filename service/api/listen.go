package api

import (
	"Rpc/decode"
)

func (s *Services) Listen() {

	for {
		bytes := make([]byte, 1024)
		read, err := s.con.Read(bytes)
		if err != nil {
			return
		}
		c := decode.Decode(bytes[0 : read+1])
		switch c.OpenCode {
		case 2:
			go s.Call(c, s.con)
		case 3:
			//ping message  δΈεε€η
		}
	}
}
