package main

import (
	"fmt"
	"sirendaou.com/duserver/apns_server"
)

func main() {
	apns_server.StartServer()
	fmt.Println("start apns_server")
}
