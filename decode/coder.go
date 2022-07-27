package decode

type CMS struct {
	// c to m 服务端发送给master 的信息
	OpenCode int64                  `msg:"open_code"`
	Id       int64                  `msg:"id"`
	Message  map[string]interface{} `msg:"message"`
}

func Encode(m CMS) ([]byte, error) {
	return m.MarshalMsg(nil)
}
func Decode(b []byte) (c CMS) {
	_, err := c.UnmarshalMsg(b)
	if err != nil {
		return CMS{}
	}
	return
}
