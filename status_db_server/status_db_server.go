package status_db_server

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"sirendaou.com/duserver/common"
)

var (
	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
	g_log_file  = flag.String("log_file", "loadreg2es.log", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")

	g_MongodbAddr  = flag.String("mongodb_addr", "", "mongodb  server addr")
	g_MsgRedisAddr = flag.String("msg_redis_addr", "", "msg cache redis server addr")

	g_nsqdTCPAddrs   = flag.String("nsq_addr", "", "nsq Server address (transient)")
	g_Conn2StatTopic = flag.String("conn2stat_topic", "t-conn2statcenter", "the name of connect to status center")
)

func SendMsg2ConnMq(tail *common.InnerPkgTail, msgItem *common.UserMsgItem, w *nsq.Producer, Cmd uint16) {
	ret, head, jsonBody, _ := common.DecPkgInnerBody([]byte(msgItem.Content))
	if ret != 0 {
		logger.Error("DecPkgInnerBody fail ret ", ret)
		return
	}

	head.Cmd = Cmd
	respBuf := new(bytes.Buffer)
	binary.Write(respBuf, binary.BigEndian, head)
	binary.Write(respBuf, binary.BigEndian, jsonBody)
	binary.Write(respBuf, binary.BigEndian, *tail)

	topic := fmt.Sprintf("conn_%s_%d", common.IntToIP(tail.ConnIP), tail.ConnPort)

	//resp send to conn
	err := w.Publish(topic, respBuf.Bytes())
	if err != nil {
		logger.Error("----------------Publish ", topic, "err:", err.Error())
	} else {
		logger.Debug("----------------Publish ", topic, " success")
	}

	return
}

func ProcFindMsg2MongoDB(loginPkgCh chan common.InnerPkgTail, msgRedisMgr *common.RedisManager) {
	session, err := mgo.Dial(*g_MongodbAddr)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	//Optional. Switch the session to a monotonic behavior
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("imsdk-msg").C("msg")

	addrs := strings.Split(*g_nsqdTCPAddrs, ",")
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

	msgItem := common.UserMsgItem{}

	count := 0
	sendCount := 0

	for {
		select {
		case tail := <-loginPkgCh:
			if tail.ConnIP == 0 || tail.ConnPort == 0 || tail.Sid == 0 || tail.ToUid == 0 {
				logger.Debug("logout Uid  not triger db msg", tail.FromUid)
				continue
			}

			if tail.ToUid != 2 {
				logger.Debug("Uid not  hello cmd not triger db msg", tail.FromUid)
				continue
			}

			//get chat msg
			iter := c.Find(bson.M{"touid": tail.FromUid}).Sort("msgid").Limit(6).Iter()
			count = 0

			for iter.Next(&msgItem) {
				logger.Info("charmsgItem ", count, "touid:", msgItem.ToUid, "msgid:", msgItem.MsgId)
				SendMsg2ConnMq(&tail, &msgItem, producers[count%producerCount], common.DU_PUSH_CMD_IM_USER_MSG)
				count++
				sendCount++
				time.Sleep(20 * time.Millisecond)
			}

			if err := iter.Close(); err != nil {
				logger.Error("iter close err:", err)
			}

			if count >= 6 {
				logger.Debug("sendCount:", count)
				continue
			}

			// get team msg
			key := fmt.Sprintf("%s%d", common.SET_TEAMMSGID, tail.FromUid)
			msgIds := msgRedisMgr.RedisZRange2(key, 6)
			for _, strMsgId := range msgIds {
				nMsgId, err := strconv.ParseUint(strMsgId, 10, 64)
				if err != nil {
					logger.Info("strmsgid :", strMsgId, " is error")
					msgRedisMgr.RedisZRem(key, strMsgId)
					continue
				}

				msgKey := fmt.Sprintf("%s%s", common.KEY_TEAMMSGBUF, strMsgId)
				val, errcode := msgRedisMgr.RedisGet(msgKey)
				if errcode != 0 {
					logger.Info("get :", msgKey, " errcode:", errcode)
					continue
				} else if len(val) <= 24 {
					logger.Info("get :", msgKey, " is empty")
					msgRedisMgr.RedisZRem(key, strMsgId)
					continue
				}

				msgItem := common.UserMsgItem{
					MsgId:    nMsgId,
					FromUid:  0,
					ToUid:    tail.FromUid,
					Type:     0,
					Content:  val,
					SendTime: 0,
					ApnsText: "",
					ExtraData:"",
				}

				logger.Info("msgItem ", count, msgItem.ToUid, msgItem.MsgId)

				SendMsg2ConnMq(&tail, &msgItem, producers[count%producerCount], common.DU_PUSH_CMD_IM_TEAM_MSG)

				count++
				sendCount++
				time.Sleep(20 * time.Millisecond)
			}

			if sendCount == 0 {
				//have no msgdb ,clear the flag
				//				common.MsgDBFlagClear(tail.FromUid)
				logger.Info("login Uid have no msg ,clear flag", tail.FromUid)
			}

		}
	}
}

type Handler struct {
	count      int
	loginPkgCh chan common.InnerPkgTail
	redisMgr *common.RedisManager
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	h.count++

	reader := bytes.NewReader(message.Body)
	tail := common.InnerPkgTail{}
	err := binary.Read(reader, binary.BigEndian, &tail)
	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return err
	}

	logger.Debug("tail:", tail.ToString())

	//check have offline msg
//	state := common.MsgDBFlagGet(tail.FromUid)
	state, err := h.redisMgr.RedisMsgCacheGet(tail.FromUid)
	if err != nil {
		return err
	}

	if state == "read" {
		logger.Debug(tail.FromUid, " have no db msg")
	} else {
		logger.Debug(tail.FromUid, " maybe have db msg")
		h.loginPkgCh <- tail
		h.redisMgr.RedisMsgCacheSet(tail.FromUid, "read")
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

	if len(*g_nsqdTCPAddrs) < 1 {
		flag.PrintDefaults()
		return
	}

	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0

	consumer, err := nsq.NewConsumer(*g_Conn2StatTopic, "msgdb-channel", nsq.NewConfig())
	if err != nil {
		logger.Error("NewConsumer ", *g_Conn2StatTopic, " error:", err.Error())
		return
	}

	defer consumer.Stop()

	consumer.SetLogger(nil, nsq.LogLevelInfo)

	ch := make(chan int)
	loginPkgCh := make(chan common.InnerPkgTail, 1000)

	msgRedis := common.NewRedisManager(*g_MsgRedisAddr)
	if msgRedis == nil {
		logger.Info("connect user redis ", *g_MsgRedisAddr, "fail")
		return
	}

	i := 0
	for i < 5 {
		go ProcFindMsg2MongoDB(loginPkgCh, msgRedis)
		i++
	}

	handler := &Handler{0, loginPkgCh, msgRedis}
	consumer.AddHandler(handler)

	addrs := strings.Split(*g_nsqdTCPAddrs, ",")
	err = consumer.ConnectToNSQDs(addrs)

	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	select {
	case ret := <-ch:
		fmt.Println("quit:", ret)
	}

	return
}
