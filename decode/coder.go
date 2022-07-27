package decode

// 信息交互方式  采用messagepackege 进行数据交互
// 据说比json 更好 有更好的压缩 编解码也更快
type CMS struct {
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
