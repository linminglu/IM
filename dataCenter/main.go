package main

import (
	"log"
	"time"

	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"

	"IM/common/rpc"
)

var (
	g_cache cache.Cache = nil
)

func registerService(obj interface{}) {
	rpc.Register(obj)
}

func main() {
	log.SetFlags(log.Lshortfile)

	redisCache, err := cache.NewCache("redis", `{"conn":":6379"}`)
	if err != nil {
		panic(err)
	}
	g_cache = redisCache

	go rpc.Serve(":9000")

	time.Sleep(time.Second)
	client, err := rpc.Dail(":9000")
	if err != nil {
		log.Println(err)
		return
	}

	arg := rpc.NewArgType()
	arg.Args["userId"] = 1
	arg.Args["status"] = 0
	if err := client.Call("UserService.SetStatus", arg, nil); err != nil {
		log.Println(err)
		return
	}

	reply := rpc.NewReplyType()

	if err := client.Call("UserService.GetStatus", arg, &reply); err != nil {
		log.Println(err)
		return
	}
	log.Println(reply.Reply)

	time.Sleep(time.Second)
}
