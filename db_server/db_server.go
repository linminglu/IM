package db

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"

	"sirendaou.com/duserver/common"
)

var (
	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")
	//	g_isdaemon       = flag.Int("isdaemon", 0, "is run server as daemon 0-no 1-s)")

	g_log_path  = flag.String("log_path", "", "the log file path")
	g_log_file  = flag.String("log_file", "", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_mysql_host = flag.String("mysql_host", "", "mysql host")
	g_mysql_db   = flag.String("mysql_db", "", "mysql db name")
	g_mysql_user = flag.String("mysql_user", "", "mysql user")
	g_mysql_pwd  = flag.String("mysql_pwd", "", "mysql passwd")

	//	lifetime         = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")

	g_RedisAddr      = flag.String("redis_addr", "", "redis mq server addr")
	g_LocRedisAddr = flag.String("location_redis_addr", "", "local cache redis  server addr")
	g_AppRedisAddr   = flag.String("app_redis_addr", "", "app redis server addr")

	//	g_WBRedisAddr    = flag.String("wbredisaddr", "", "whilte black list redis server addr")
	g_MongodbAddr = flag.String("mongodb_addr", "", "mongodb  server addr")

	g_nsqd_addrs = flag.String("nsq_addr", "", "nsq Server address (transient)")

	g_Conn2DbTopic = flag.String("conn2db_topic", "t-conn2db", "the name of connect to db topic")
	g_Db2MsgTopic  = flag.String("db2msg_topic", "t-db2msgcenter", "the name of db to msg center")

	//	g_RbAddr         = flag.String("addr", "", "rabbitmq server addr")

//	uri          = flag.String("uri", "amqp://du:MQ_du@182.254.178.103:5672/", "AMQP URI")
//	exchange     = flag.String("product-exchange", "ex-test", "Durable, non-auto-deleted AMQP exchange name")
//	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
)

func SendMsg2MsgCenter(msgBodyCh chan []byte) int {
	addrs := strings.Split(*g_nsqd_addrs, ",")

	producerCount := len(addrs)
	producers := make([]*nsq.Producer, producerCount)

	var err error
	for i, add := range addrs {
		producers[i], err = nsq.NewProducer(add, nsq.NewConfig())
		if err != nil {
			logger.Error("NewProducer  ", add, " err:", err)
			return -1
		}
	}

	defer func() {
		for _, producer := range producers {
			producer.Stop()
		}
	}()

	count := uint64(0)
	for {
		select {
		case msg := <-msgBodyCh:
			count++
			err = producers[int(count)%producerCount].Publish(*g_Db2MsgTopic, msg)
			if err != nil {
				logger.Error("nsq write resp to ", *g_Db2MsgTopic, "err:", err.Error())
			} else {
				logger.Debug("-----------SendMsg2MsgCenter nsq write resp to ", *g_Db2MsgTopic, " success")
			}
		}
	}

	return 0
}

type UserLocation struct {
	Uid      uint32
	Loc      []float64
	SendTime uint32
}

type DBHandler struct {
	count      int
	ch         chan int
	chatMsgCh  chan []byte
	producer   []*nsq.Producer
	UserRedis  *common.RedisManager
	AppRedis   *common.RedisManager
	LocRedis *common.RedisManager
}

func (h *DBHandler) HandleMessage(message *nsq.Message) error {
	h.count++

	p := bytes.NewReader(message.Body)
	head := common.PkgHead{}
	err := binary.Read(p, binary.BigEndian, &head)

	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return err
	}

	logger.Debug("HandleMessage PkgHead ", head)

	if int(head.PkgLen)+common.SIZEOF_INNERTAIL != len(message.Body) {
		logger.Error("pkghead len(", head.PkgLen, ") , pkg len(", len(message.Body), ") is error")
		return nil
	}

	jsonStr := make([]byte, head.PkgLen-common.SIZEOF_PKGHEAD)
	err = binary.Read(p, binary.BigEndian, &jsonStr)

	if err != nil {
		logger.Error("read req json fail:", err.Error())
		return err
	}

	logger.Info("req json:", string(jsonStr[:]))

	tail := common.InnerPkgTail{}
	err = binary.Read(p, binary.BigEndian, &tail)
	if err != nil {
		logger.Error("read req innertail fail:", err.Error())
		return err
	}

	logger.Info("req head:", head.ToString())
	logger.Info("req innertail:", tail.ToString())
	logger.Info("req json:", string(jsonStr[:]))

	var resp []byte
	var retCode uint32 = 0

	switch head.Cmd {
	case common.DU_CMD_USER_REGISTER:
		resp, retCode = h.Register(head, jsonStr, &tail)
//	case common.DU_CMD_USER_PURE_REGISTER:
//		resp, retcode = h.PureReg(head, jsonstr, &tail)
	case common.DU_CMD_USER_LOGIN:
		resp, retCode = h.Login(head, jsonStr, &tail)
	case common.DU_CMD_USER_TOKEN_LOGIN:
		resp, retCode = h.LoginWithToken(head, jsonStr, &tail)
	case common.DU_CMD_USER_SET_DEVICE_TOKEN:
		resp, retCode = h.SetDeviceToken(head, jsonStr, &tail)
	case common.DU_CMD_USER_SET_MY_INFO:
		resp, retCode = h.SetUserInfo(head, jsonStr, &tail)
	case common.DU_CMD_USER_GET_INFO:
		resp, retCode = h.GetUserInfo(head, jsonStr, &tail)
	case common.DU_CMD_USER_GET_UID:
		resp, retCode = h.GetUid(head, jsonStr, &tail)
	case common.DU_CMD_AROUND_QUERY:
		resp, retCode = h.RequestUserLocation(head, jsonStr, &tail)
	case common.DU_CMD_USER_RESET_PWD:
		resp, retCode = h.ResetPwd(head, jsonStr, &tail)
	case common.DU_CMD_USER_BIND_PHONE:
		resp, retCode = h.BindPhone(head, jsonStr, &tail)
	case common.DU_CMD_USER_RETRIEVE_PWD:
		resp, retCode = h.RetrievePwd(head, jsonStr, &tail)
		//	case common.DU_CMD_USER_GET_APP_INFO:
		//		resp, retCode = h.AppCfg(head, jsonStr, &tail)
		//	case common.DU_PUSH_CMD_IM_REPORT_MSG:
		//		resp, retCode = h.Report(head, jsonStr, &tail)
		//	case common.DU_CMD_USER_GET_SETUP_ID:
		//		resp, retCode = h.GetSetupId(head, jsonStr, &tail)
		//	case common.DU_CMD_USER_SET_SETUP_ID:
		//		resp, retCode = h.SetSetupId(head, jsonStr, &tail)
		//	case common.DU_CMD_GET_CSID_LIST:
		//		resp, retCode = h.GetCSIdList(head, jsonStr, &tail)
		//	case common.DU_CMD_GET_CSINFO_LIST:
		//		resp, retCode = h.GetCSList(head, jsonStr, &tail)

	default:
		logger.Debug("invalid cmd:", head.Cmd)
		return nil
	}

	respBuf := new(bytes.Buffer)

	if retCode != 0 {
		head.PkgLen = uint16(common.SIZEOF_PKGHEAD)
		head.Sid = uint32(retCode)
	} else {
		head.PkgLen = uint16(common.SIZEOF_PKGHEAD + len(resp))
		head.Sid = uint32(0)
		head.Uid = tail.FromUid
	}

	binary.Write(respBuf, binary.BigEndian, head)

	if retCode == 0 && len(resp) > 0 {
		binary.Write(respBuf, binary.BigEndian, resp)
	}

	binary.Write(respBuf, binary.BigEndian, tail)

	topic := fmt.Sprintf("conn_%s_%d", common.IntToIP(tail.ConnIP), tail.ConnPort)

	//resp send to conn
	if len(topic) > 0 {
		err := h.producer[h.count%(len(h.producer))].Publish(topic, respBuf.Bytes())
		if err != nil {
			logger.Error("nsq write resp to ", topic, "err:", err.Error())
		} else {
			logger.Debug("nsq write resp to ", topic, " success")
		}
	}

	return nil
}

func Handle() {
	var err error

	if common.MysqlInit(*g_mysql_host, *g_mysql_db, *g_mysql_user, *g_mysql_pwd) != 0 {
		logger.Error("mysql init ", *g_mysql_host, *g_mysql_db, *g_mysql_user, "fail")
		return
	}

	if len(*g_nsqd_addrs) < 1 {
		flag.PrintDefaults()
		return
	}

	userRedis := common.NewRedisManager(*g_RedisAddr)
	if userRedis == nil {
		logger.Info("connect user redis ", *g_RedisAddr, "fail")
		return
	}

	logger.Info("connect user redis ", *g_RedisAddr, " success")

	appRedis := common.NewRedisManager(*g_AppRedisAddr)
	if appRedis == nil {
		logger.Info("connect user redis ", *g_AppRedisAddr, "fail")
		return
	}

	logger.Info("connect app redis ", *g_AppRedisAddr, " success")

	localRedis := common.NewRedisManager(*g_LocRedisAddr)
	if localRedis == nil {
		logger.Info("connect user redis ", *g_LocRedisAddr, "fail")
		return
	}

	logger.Info("connect local redis ", *g_LocRedisAddr, " success")

	if common.MongoInit(*g_MongodbAddr) != 0 {
		logger.Info("connect mongo ", *g_MongodbAddr, "fail")
		return
	}

	logger.Info("connect mongo ", *g_MongodbAddr, " success")

	addrs := strings.Split(*g_nsqd_addrs, ",")
	producerCount := len(addrs)
	producers := make([]*nsq.Producer, producerCount)

	for i, add := range addrs {
		config := nsq.NewConfig()
		config.DefaultRequeueDelay = 0
		producers[i], err = nsq.NewProducer(add, config)

		if err != nil {
			logger.Error("NewProducer  ", add, " err:", err)
			return
		}
	}

	defer func() {
		for _, producer := range producers {
			producer.Stop()
		}
	}()

	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0
	consumer, err := nsq.NewConsumer(*g_Conn2DbTopic, "db-channel", config)

	if err != nil {
		logger.Error("NewConsumer ", *g_Conn2DbTopic, " error:", err.Error())
		return
	}

	defer consumer.Stop()

	consumer.SetLogger(nil, nsq.LogLevelInfo)

	ch := make(chan int)

	chatMsgCh := make(chan []byte, 1000)

	go SendMsg2MsgCenter(chatMsgCh)

	//	if comm.RbmqInit(*uri) != 0 {
	//		logger.Error("RbmqInit  fail ")
	//		return
	//	}
	//
	//	if comm.RbmqDeclareExchange(*exchange, *exchangeType) != 0 {
	//		return
	//	}

	handler := &DBHandler{0, ch, chatMsgCh, producers, userRedis, appRedis, localRedis}

	consumer.AddHandler(handler)

	addrs = strings.Split(*g_nsqd_addrs, ",")
	err = consumer.ConnectToNSQDs(addrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	select {}

	return
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

	fmt.Println("StartServer")
	logger.Info("StartServer")

	//	for i := 0; i < 8; i++ {
	for i := 0; i < 2; i++ {
		go Handle()
		time.Sleep(time.Second * 2)
	}

	select {}
}
