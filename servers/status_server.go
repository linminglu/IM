package main

import (
	"fmt"
	"sirendaou.com/duserver/status_server"
)

func main() {
	status_center.StartServer()
	fmt.Println("start db server")
}
