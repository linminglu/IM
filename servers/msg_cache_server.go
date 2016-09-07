package main

import (
	"fmt"
	"sirendaou.com/duserver/msg_cache_server"
)

func main() {
	msg_cache_server.StartServer()
	fmt.Println("start msg_cache_server")
}
