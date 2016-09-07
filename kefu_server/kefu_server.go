package kefu

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"
	"github.com/rakyll/globalconf"

	"sirendaou.com/duserver/common"
)

var (
	g_ListenAddr = flag.String("listenaddr", "0.0.0.0:8080", "redis mq server addr")
	g_cpu_num    = flag.Int("cpu_num", 4, "the num of cpu")

	g_log_path    = flag.String("log", "/tmp", "the log file path")
	g_log_file    = flag.String("log_file", "loadreg2es.log", "the log file path")
	g_log_level   = flag.Int("log_level", 2, "the log level 1-debug 2-info(default) 3-WARN 4-error 5-FATAL 6-off")
	g_svrlog_path = flag.String("svrlog", "/tmp/dlsvr.log", "the framwork log file path")

	//	g_connaddr    = flag.String("connaddr", "", "the du conn addr")

	g_mysql_host = flag.String("mysql_host", "", "mysql host")
	g_mysql_db   = flag.String("mysql_db", "", "mysql db name")
	g_mysql_user = flag.String("mysql_user", "", "mysql user")
	g_mysql_pwd  = flag.String("mysql_pwd", "", "mysql passwd")
	//	g_mysql_host2   = flag.String("mysql_host_portal", "", "portal mysql host")
	//	g_mysql_db2     = flag.String("mysql_db_portal", "", "portal mysql db name")
	//	g_mysql_user2   = flag.String("mysql_user_portal", "", "portal mysql user")
	//	g_mysql_pwd2 = flag.String("mysql_pwd_portal", "", "portal mysql passwd")

	g_RedisAddr    = flag.String("redis_addr", "", "redis mq server addr")
	g_MsgRedisAddr = flag.String("msg_redis_addr", "", "kefu msg redis mq server addr")
	g_AppRedisAddr = flag.String("app_redis_addr", "", "app redis server addr")
	g_nsqd_addrs   = flag.String("nsq_addr", "", "nsq Server address (transient)")
	g_Db2MsgTopic  = flag.String("db2msg_topic", "t-db2msgcenter", "the name of db to msg center")

	g_DownloadUrl = flag.String("downloadurl", "", "7niu download url")
	g_MongodbAddr = flag.String("mongodb_addr", "", "mongodb  server addr")
)

type Handler struct {
	UserRedis *common.RedisManager
	//	AppRedis  *common.RedisManager
	MsgRedis *common.RedisManager
	Producer *nsq.Producer
	Session  map[string]string
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

	logger.SetConsole(false)
	logger.SetRollingDaily(*g_log_path, *g_log_file)
	logger.SetLevel(logger.LEVEL(*g_log_level))

	logger.Info("start:")

	logfile, err := os.OpenFile(*g_svrlog_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}

	defer logfile.Close()

	svrLog := log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

	producer, err := nsq.NewProducer(*g_nsqd_addrs, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	if common.MysqlInit(*g_mysql_host, *g_mysql_db, *g_mysql_user, *g_mysql_pwd) != 0 {
		logger.Error("mysql init ", *g_mysql_host, *g_mysql_db, *g_mysql_user, "fail")
		return
	}

	//	AppMap := make(map[string]int)
	//	common.DBLoadAppkey(AppMap)

	if len(*g_nsqd_addrs) < 1 {
		flag.PrintDefaults()
		return
	}

	userRedis := common.NewRedisManager(*g_RedisAddr)
	if userRedis == nil {
		logger.Info("connect user redis ", *g_RedisAddr, "fail")
		return
	}
	logger.Info("connect  user redis ", *g_RedisAddr, " success")

	//	appredis := common.NewRedisManager(*g_AppRedisAddr)
	//	if appredis == nil {
	//		logger.Info("connect user redis ", *g_AppRedisAddr, "fail")
	//		return
	//	}
	//	logger.Info("connect  app redis ", *g_AppRedisAddr, " success")

	msgRedis := common.NewRedisManager(*g_MsgRedisAddr)
	if msgRedis == nil {
		logger.Info("connect user redis ", *g_MsgRedisAddr, "fail")
		return
	}
	logger.Info("connect  msg redis ", *g_MsgRedisAddr, " success")

	if common.MongoInit(*g_MongodbAddr) != 0 {
		logger.Info("connect mongo ", *g_MongodbAddr, "fail")
		return
	}

	logger.Info("connect mongo ", *g_MongodbAddr, " success")

	cookie := map[string]string{}

	//	handler := Handler{userRedis, appredis, msgRedis, producer, cookie}
	handler := Handler{userRedis, msgRedis, producer, cookie}

	web.Post("(/cs/init)", handler.Init)
	web.Get("(/cs/init)", handler.Init)
	web.Post("(/cs/setup)", handler.Init)
	web.Get("(/cs/setup)", handler.Init)

	web.Post("(/cs/account/register)", handler.Register)
	web.Get("(/cs/account/register)", handler.Register)
	web.Post("(/cs/account/enable)", handler.Enable)
	web.Get("(/cs/account/enable)", handler.Enable)
	web.Post("(/cs/account/delete)", handler.Delete)
	web.Get("(/cs/account/delete)", handler.Delete)
	web.Post("(/cs/account/login)", handler.Login)
	web.Get("(/cs/account/login)", handler.Login)
	web.Post("(/cs/account/logout)", handler.LogOut)
	web.Get("(/cs/account/logout)", handler.LogOut)
	web.Post("(/cs/account/list)", handler.List)
	web.Get("(/cs/account/list)", handler.List)
	web.Post("(/cs/account/current)", handler.Current)
	web.Get("(/cs/account/current)", handler.Current)
	web.Post("(/cs/account/password)", handler.Password)
	web.Get("(/cs/account/password)", handler.Password)

	web.Post("(/cs/msg/list)", handler.MsgList)
	web.Get("(/cs/msg/history)", handler.History)
	web.Post("(/cs/msg/history)", handler.History)
	web.Get("(/cs/msg/list)", handler.MsgList)
	web.Post("(/cs/msg/send)", handler.Send)
	web.Get("(/cs/msg/send)", handler.Send)
	web.Post("(/cs/msg/recv)", handler.Recv)
	web.Get("(/cs/msg/recv)", handler.Recv)

	web.Get("(/cs/image/download)", handler.ImageDownload)

	web.SetLogger(svrLog)
	web.Run(*g_ListenAddr)
}
