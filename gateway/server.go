package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"

	"IM/common/json"
	sync_ "IM/common/sync"
)

type Gateway struct {
	config      config.Configer
	clientMap   map[int]*Client
	clientMutex *sync.Mutex
	waitGroup   *sync_.WaitGroup
}

var (
	g_gateway = &Gateway{
		clientMap:   make(map[int]*Client),
		clientMutex: &sync.Mutex{},
		waitGroup:   sync_.NewWaitGroup(),
	}
)

func StartServer() {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	ip := g_gateway.config.DefaultString("serverIp", "")
	port := g_gateway.config.DefaultString("serverPort", "9100")

	addr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		logs.Debug(err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logs.Debug(err)
		return
	}

	logs.Debug("StartServer listening on:", listener.Addr())

	var sid int = 10000
	for {
		listener.SetDeadline(time.Now().Add(time.Millisecond * 20))
		conn, err := listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			logs.Debug("Network Error:", err)
			return
		}
		logs.Debug("Client Connect, addr:", conn.RemoteAddr().String())
		// 建立客户端
		client := &Client{
			Conn:      conn,
			Sid:       sid,
			ReqMsgCh:  make(chan *json.Json, 100),
			RespMsgCh: make(chan *json.Json, 100),
			waitGroup: sync_.NewWaitGroup(),
		}

		// 处理客户端请求

		g_gateway.clientMutex.Lock()
		g_gateway.clientMap[sid] = client
		g_gateway.clientMutex.Unlock()

		sid++ // increase sid

		go client.Serve()
	}
	return
}

func loadConfig() error {
	configFile := "./config.ini"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config, err := config.NewConfig("ini", configFile)
	if err != nil {
		log.Fatal("read config:", err)
		return err
	}

	g_gateway.config = config
	return nil
}
