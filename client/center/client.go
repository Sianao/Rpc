package center

import (
	"Rpc/client/logic"
	"Rpc/decode"
	"log"
	"net"
	"strconv"
	"strings"
)

// 定义客户端信息

type Client struct {
	Con     net.Conn //网路连接
	Service string   // 调用
	CMs     decode.CMS
}

// NewClient  新建一个客户端
func NewClient(addr ...string) (Client, error) {
	add := ":4563"
	if len(addr) != 0 {
		add = addr[0]
	}
	con, err := net.Dial("tcp", add)
	if err != nil {
		return Client{}, err
	}
	return Client{Con: con}, nil
}

// Call  表明调用的服务
func (c *Client) Call(service string) *Client {
	// 表明这是一个调用请求
	c.Service = service
	return c
}

// FUN 表明调用的方法 并等
func (c *Client) FUN(me string, param ...interface{}) ([]interface{}, error) {
	//service 包含其调用的方法 调用的
	c.CMs.OpenCode = 2
	// 根据发出请求的端口 生成id 构造请求
	id, _ := strconv.Atoi(strings.Split(c.Con.LocalAddr().String(), ":")[1])
	c.CMs.Id = int64(id)
	c.CMs.Message = make(map[string]interface{})
	c.CMs.Message["Service"] = c.Service
	c.CMs.Message["Method"] = me
	c.CMs.Message["Param"] = param
	msg, err := logic.Call(c.Con, c.CMs)
	if err != nil {
		log.Println(err)
	}
	return msg, err
}
