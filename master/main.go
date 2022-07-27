package main

import (
	"Rpc/master/cmd"
	"flag"
)

var b = flag.String("p", "4563", "程序运行端口")

func main() {
	flag.Parse()
	cmd.Start(b)

}
