package main

import (
	"fmt"
	"sirendaou.com/duserver/msg_server"
)

func main() {
	msg_server.StartServer()
	fmt.Println("start msg server")
}
