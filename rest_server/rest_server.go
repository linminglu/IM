package rest_server

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"

	"github.com/go-martini/martini"
	"sirendaou.com/duserver/common"
)

var (
	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")

	g_log_path      = flag.String("log_path", "", "the log file path")
	g_log_file      = flag.String("log_file", "", "the log file path")
	g_log_level     = flag.Int("log_level", 2, "the log level 1-debug 2-info(default) 3-WARN 4-error 5-FATAL 6-off")
	g_rest_log_path = flag.String("rest_log_path", "", "the framwork log file path")

	g_mysql_host = flag.String("mysql_host", "", "mysql host")
	g_mysql_db   = flag.String("mysql_db", "", "mysql db name")
	g_mysql_user = flag.String("mysql_user", "", "mysql user")
	g_mysql_pwd  = flag.String("mysql_pwd", "", "mysql passwd")

	g_RedisAddr    = flag.String("redis_addr", "", "redis mq server addr")
	g_MsgRedisAddr = flag.String("msg_redis_addr", "", "team msg cache redis server addr")

	g_nsqd_addrs  = flag.String("nsq_addr", "", "nsq Server address (transient)")
	g_Db2MsgTopic = flag.String("db2msg_topic", "t-db2msgcenter", "the name of db to msg center")

	g_MongodbAddr = flag.String("mongodb_addr", "", "mongodb  server addr")
)

type RestResp struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}

type Handler struct {
	redisMgr *common.RedisManager
	producer *nsq.Producer
}

func StartServer() {
	if len(os.Args) < 2 {
		fmt.Println("please set conf file ")
		return
	}

	conf, err := globalconf.NewWithOptions(&globalconf.Options{
		Filename: os.Args[1],
	})

	if err != nil {
		fmt.Print("NewWithFilename ", os.Args[1], " fail :", err)
		return
	}

	conf.ParseAll()

	runtime.GOMAXPROCS(*g_cpu_num)

	//	logger.SetConsole(false)
	logger.SetRollingDaily(*g_log_path, *g_log_file)
	logger.SetLevel(logger.LEVEL(*g_log_level))

	logfile, err := os.OpenFile(*g_rest_log_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()

	if common.MysqlInit(*g_mysql_host, *g_mysql_db, *g_mysql_user, *g_mysql_pwd) != 0 {
		logger.Error("mysql init ", *g_mysql_host, *g_mysql_db, *g_mysql_user, "fail")
		return
	}

	producer, err := nsq.NewProducer(*g_nsqd_addrs, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	if common.MysqlInit(*g_mysql_host, *g_mysql_db, *g_mysql_user, *g_mysql_pwd) != 0 {
		logger.Error("mysql init ", *g_mysql_host, *g_mysql_db, *g_mysql_user, "fail")
		return
	}

	if len(*g_nsqd_addrs) < 1 {
		flag.PrintDefaults()
		return
	}

	redisMgr := common.NewRedisManager(*g_RedisAddr)
	if redisMgr == nil {
		logger.Info("connect user redis ", *g_RedisAddr, "fail")
		return
	}
	logger.Info("connect  user redis ", *g_RedisAddr, " success")

	if TeamMsgRedisInit(*g_MsgRedisAddr) != 0 {
		logger.Info("connect msg cache redis ", *g_MsgRedisAddr, "fail")
		return
	}
	logger.Info("connect msg cache redis ", *g_MsgRedisAddr, " success")

	if common.MongoInit(*g_MongodbAddr) != 0 {
		logger.Info("connect mongo ", *g_MongodbAddr, "fail")
		return
	}

	logger.Info("connect mongo ", *g_MongodbAddr, " success")

	handler := Handler{redisMgr, producer}

	m := martini.Classic()
	m.Post("(/rest/message/send)", handler.SendMessage)
	m.Post("(/rest/users)", handler.Users)
	m.Get("(/rest/chatroom/list)", handler.ListSysTeam)
	m.Post("(/rest/chatroom/list)", handler.ListSysTeam)
	m.Post("(/rest/chatroom/add)", handler.CreateSysTeam)
	m.Post("(/rest/offline_msg/list)", handler.OffLineMsg)
	m.Post("(/rest/setupid/get)", handler.GetSetupID)
	m.RunPort("3001")
}
