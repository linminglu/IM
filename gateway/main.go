package main

func main() {

	if err := loadConfig(); err != nil {
		println("load config failed, exited!!!")
	}

	StartServer()
}
