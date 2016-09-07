package main

import (
	"fmt"
	"sirendaou.com/duserver/rest_server"
)

func main() {
	rest_server.StartServer()
	fmt.Println("start rest server")
}
