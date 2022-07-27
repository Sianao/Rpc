package main

import (
	"Rpc/client/center"
	"fmt"
)

func main() {

	client, err := center.NewClient(":4563")
	if err != nil {
		panic(err)
	}
	re, err := client.Call("MS").FUN("Add")
	if err != nil {
		return
	}
	fmt.Println(re)
}
