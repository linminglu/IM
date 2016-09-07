package main

import (
	"fmt"
	"sirendaou.com/duserver/file_server"
)

func main() {
	file_server.StartServer()
	fmt.Println("start file_server")
}
