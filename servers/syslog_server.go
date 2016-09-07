package main

import (
	"fmt"
	"sirendaou.com/duserver/syslog_server"
)

func main() {
	syslog_server.StartServer()
	fmt.Println("start syslog server")
}
