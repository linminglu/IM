package main

import (
	"fmt"
	"sirendaou.com/duserver/status_db_server"
)

func main() {
	status_db_server.StartServer()
	fmt.Println("start db server")
}
