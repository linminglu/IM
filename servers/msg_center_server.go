package main

import (
	"fmt"
	"sirendaou.com/duserver/msg_center"
)

func main() {
	msg_center.StartServer()
	fmt.Println("start msg_center_server")
}
