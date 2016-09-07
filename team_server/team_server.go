package team

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"encoding/binary"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"sirendaou.com/duserver/common"
)

var (
//	g_RbAddr = flag.String("addr", "", "rabbitmq server addr")

	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
	g_log_file  = flag.String("log_file", "loadreg2es.log", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")

//	g_isdaemon   = flag.Int("isdaemon", 0, "is run server as daemon 0-no 1-s)")

	g_mysql_host = flag.String("mysql_host", "", "mysql host")
	g_mysql_db   = flag.String("mysql_db", "", "mysql db name")
	g_mysql_user = flag.String("mysql_user", "", "mysql user")
	g_mysql_pwd  = flag.String("mysql_pwd", "", "mysql passwd")

	g_RedisAddr     = flag.String("redis_addr", "", "team info redis server addr")
	g_MsgRedisAddr  = flag.String("msg_redis_addr", "", "team msg cache redis server addr")
	g_UserRedisAddr = flag.String("user_redis_addr", "", "user info redis server addr")

	g_MongodbAddr   = flag.String("mongodb_addr", "", "mongodb  server addr")

	g_Conn2TeamTopic = flag.String("conn2team_topic", "t-conn2team", "the name of connect to db topic")
	g_DbMsgTopic  = flag.String("db2msg_topic", "t-db2msgcenter", "the name of db to msg center")
	g_nsqd_addrs   = flag.String("nsq_addr", "", "nsq Server address (transient)")

//	lifetime = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
)

func ProcSaveMsg2MongoDB(msgSaveCh chan *common.TeamMsgSaveItem) {
	session, err := mgo.Dial(*g_MongodbAddr)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("team").C("msg")

	for {
		select {
		case msg := <-msgSaveCh:
			logger.Info("mongdb teamMsg to Uid", msg.TouId, " msgid ", msg.MsgId, " opt ", msg.Opt)
			if msg.Opt == 1 {
				err = c.Insert(msg)
			} else if msg.Opt == 0 {
				//strtemp := fmt.Sprintf(`{"Msgid":%d, "Touid":%d}`, msg.Msgid, msg.Touid)
				bs := bson.M{"msgid": msg.MsgId, "touid": msg.TouId}
				logger.Debug("bson:", bs)
				err = c.Remove(bs)
			} else {
				err = fmt.Errorf(" error  Opt", msg.Opt)
			}

			if err != nil {
				logger.Error("msg msgid:", msg.MsgId, "remove mongodb fail:", err.Error())
			} else {
				logger.Info("msg msgid:", msg.MsgId, "touid", msg.TouId, " remove mongodb success")
			}
		}
	}
}

func SendMsg2MsgCenter(msgBodyCh chan []byte) int {
	producer, err := nsq.NewProducer(*g_nsqd_addrs, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	defer producer.Stop()

	for {
		select {
		case msg := <-msgBodyCh:
			err = producer.Publish(*g_DbMsgTopic, msg)
			if err != nil {
				logger.Error("nsq write resp to ", *g_DbMsgTopic, "err:", err.Error())
			} else {
				logger.Debug("nsq write resp to ", *g_DbMsgTopic, " success")
			}
		}
	}

	return 0
}

type TeamHandler struct {
	count      int
	ch         chan int
	chatMsgCh  chan []byte
	msgSaveCh  chan *common.TeamMsgSaveItem
	producers []*nsq.Producer
	userRedisMgr *common.RedisManager
}

func (h *TeamHandler) HandleMessage(message *nsq.Message) error {
	h.count++

	p := bytes.NewReader(message.Body)
	head := common.PkgHead{}
	err := binary.Read(p, binary.BigEndian, &head)

	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return err
	}

	logger.Debug("PkgHead:", head)

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

	logger.Info("req innertail:", tail)
	logger.Info("req json:", string(jsonStr[:]))

	var resp []byte
	var retCode uint32 = 0

	switch head.Cmd {
	case common.DU_CMD_TEAM_CREATE:
		resp, retCode = h.Create(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_DELETE:
		resp, retCode = h.Delete(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_GET_INFO:
		resp, retCode = h.Query(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_GET_ALL:
		resp, retCode = h.QueryList(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_GET_SYS:
		resp, retCode = h.QuerySysList(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_SET_INFO:
		resp, retCode = h.Update(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_ADD_MEMBER:
		resp, retCode = h.AddMembers(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_REMOVE_MEMBER:
		resp, retCode = h.RemoveMember(head, jsonStr, &tail)
	case common.DU_CMD_TEAM_GET_MEMBER:
		resp, retCode = h.QueryMember(head, jsonStr, &tail)
	case common.DU_CMD_IM_SEND_TEAM_MSG:
		resp, retCode = h.SendMsg(head, jsonStr, &tail)
	case common.DU_CMD_ADD_MEMBER2WB:
		resp, retCode = h.WBAdd(head, jsonStr, &tail)
	case common.DU_CMD_DEL_MEMBER2WB:
		resp, retCode = h.WBDelete(head, jsonStr, &tail)
	case common.DU_CMD_GET_MEMBER2WB:
		resp, retCode = h.WBQuery(head, jsonStr, &tail)
	case common.DU_CMD_IM_TEAM_MSG_RECEIVED:
		resp, retCode = h.RecvedMsg(head, jsonStr, &tail)
	default:
		logger.Debug(" error Cmd :", head.Cmd)
		return nil
	}

	logger.Info("Cmd", head.Cmd, " resp:", string(resp[:]), " retcode:", retCode)

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
		err := h.producers[h.count%len(h.producers)].Publish(topic, respBuf.Bytes())
		if err != nil {
			logger.Error("nsq write resp to ", topic, "err:", err.Error())
		} else {
			logger.Debug("nsq write resp to ", topic, " success")
		}
	}

	return nil
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

	if common.MysqlInit(*g_mysql_host, *g_mysql_db, *g_mysql_user, *g_mysql_pwd) != 0 {
		logger.Error("mysql init ", *g_mysql_host, *g_mysql_db, *g_mysql_user, "fail")
		return
	}

	if len(*g_nsqd_addrs) < 1 {
		flag.PrintDefaults()
		return
	}

	if RedisInit(*g_RedisAddr) != 0 {
		logger.Info("connect redis ", *g_RedisAddr, "fail")
		return
	}
	logger.Info("connect redis ", *g_RedisAddr, " success")

	if TeamMsgRedisInit(*g_MsgRedisAddr) != 0 {
		logger.Info("connect msg cache redis ", *g_MsgRedisAddr, "fail")
		return
	}
	logger.Info("connect msg cache redis ", *g_MsgRedisAddr, " success")

	userRedis := common.NewRedisManager(*g_UserRedisAddr)
	if userRedis == nil {
		logger.Info("connect user redis ", *g_UserRedisAddr, "fail")
		return
	}

	logger.Info("connect user redis ", *g_UserRedisAddr, " success")

	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0

	addrs := strings.Split(*g_nsqd_addrs, ",")
	producerCount := len(addrs)
	producers := make([]*nsq.Producer, producerCount)

	for i, add := range addrs {
		producers[i], err = nsq.NewProducer(add, nsq.NewConfig())
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

	r, err := nsq.NewConsumer(*g_Conn2TeamTopic, "db-channel", nsq.NewConfig())
	if err != nil {
		logger.Error("NewConsumer ", *g_Conn2TeamTopic, " error:", err.Error())
		return
	}

	defer r.Stop()

	r.SetLogger(nil, nsq.LogLevelInfo)

	ch := make(chan int)

	chatMsgCh := make(chan []byte, 1000)

	go SendMsg2MsgCenter(chatMsgCh)

	msgSaveCh := make(chan *common.TeamMsgSaveItem, 1000)

	go ProcSaveMsg2MongoDB(msgSaveCh)

	handler := &TeamHandler{count:0, ch:ch, chatMsgCh:chatMsgCh, msgSaveCh:msgSaveCh, producers:producers, userRedisMgr:userRedis}
	r.AddHandler(handler)

	addrs = strings.Split(*g_nsqd_addrs, ",")
	err = r.ConnectToNSQDs(addrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	select {}

	return
}
