package msg_server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

	"sirendaou.com/duserver/common"
)

var (
	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")

	g_log_path  = flag.String("log_path", "", "the log file path")
	g_log_file  = flag.String("log_file", "", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_mysql_host = flag.String("mysql_host", "", "mysql host")
	g_mysql_db   = flag.String("mysql_db", "", "mysql db name")
	g_mysql_user = flag.String("mysql_user", "", "mysql user")
	g_mysql_pwd  = flag.String("mysql_pwd", "", "mysql passwd")

	g_RedisAddr        = flag.String("redis_addr", "", "redis mq server addr")
	g_UserRedisAddr    = flag.String("uesrredis_addr", "", "user redis mq server addr")
	g_TeamMsgRedisAddr = flag.String("msg_redis_addr", "", "msg cache redis server addr")
	//	g_KefuMsgRedisAddr = flag.String("kefumsgredis_addr", "", "kefu msg mq redis server addr")
	g_MongodbAddr = flag.String("mongodb_addr", "", "mongodb  server addr")

	g_nsqd_addrs = flag.String("nsq_addr", "", "nsq Server address (transient)")

	g_Conn2MsgSvrTopic = flag.String("conn2msg_topic", "t-conn2msg_server", "the name of connect to db topic")
	g_Db2MsgTopic      = flag.String("db2msg_topic", "t-db2msgcenter", "the name of db to msg center")
)

type UserLocation struct {
	Uid      uint32
	Loc      []float64
	SendTime uint32
}

type UserMsgReq struct {
	MsgContent string `json:"msgcontent,omitempty"`
	ToUid      uint64 `json:"touid,omitempty"`
	MsgType    uint16 `json:"msgtype,omitempty"`
	ApnsText   string `json:"apnstext,omitempty"`
	FBv        int    `json:"frombv,omitempty"`
	ExtraData  string `json:"extraData,omitempty"`
}

type Handler struct {
	count           int // 消息数量
	ch              chan int
	chatMsgCh       chan []byte // 消息通道，用于多协程通信
	producers       []*nsq.Producer
	redisMgr        *common.RedisManager
	userRedisMgr    *common.RedisManager
	teamMsgRedisMgr *common.RedisManager
	//	kefuMsgRedisMgr *common.RedisManager
}

func (h Handler) SendMsg2Mq(head common.PkgHead, msg common.UserMsgItem, tail *common.InnerPkgTail) string {
	jsonBody, _ := json.Marshal(msg)

	msgBuf := new(bytes.Buffer)

	head.PkgLen = common.SIZEOF_PKGHEAD + uint16(len(jsonBody))
	head.Seq = 0
	head.Sid = 0
	head.Cmd = common.DU_PUSH_CMD_IM_USER_MSG

	binary.Write(msgBuf, binary.BigEndian, head)
	binary.Write(msgBuf, binary.BigEndian, jsonBody)
	binary.Write(msgBuf, binary.BigEndian, tail)

	h.chatMsgCh <- msgBuf.Bytes()

	return msgBuf.String()
}

func (h *Handler) SendMsg(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req UserMsgReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	rClient := <-h.redisMgr.RedisCh
	defer func() {
		h.redisMgr.RedisCh <- rClient
	}()

	//	now := time.Now().Unix()
	//	if common.IsKefuUid(req.ToUid) {
	//		logger.Info("touid:", req.ToUid, " is kefu ")
	//
	//		key := fmt.Sprintf("cs_msg_%d", req.ToUid)
	//		val := fmt.Sprintf("%d|%s|%d", now, req.MsgContent, req.MsgType)
	//		h.kefuMsgRedisMgr.RedisLPush(key, val)
	//
	//		respStr := fmt.Sprintf(`{"msgid":%d,"sendtime":%d}`, common.GetChatMsgId(), now)
	//		return []byte(respStr), 0
	//	}

	if common.CacheCheckWBMember(rClient.Client, req.ToUid, head.Uid, 2) {
		logger.Info(head.Uid, "is in", req.ToUid, " blacklist")
		return []byte(""), uint32(common.IN_BLACKLIST)
	}

	var msg common.UserMsgItem

	msg.SendTime = uint32(time.Now().Unix())
	msg.MsgId = common.GetChatMsgId()
	msg.FromUid = head.Uid
	msg.ToUid = req.ToUid
	msg.Content = req.MsgContent
	msg.Type = req.MsgType
	msg.ApnsText = req.ApnsText
	msg.ExtraData = req.ExtraData
	msg.CmdType = common.DU_PUSH_CMD_IM_USER_MSG

	tail.ToUid = msg.ToUid
	tail.MsgId = msg.MsgId
	tail.ToUid = req.ToUid

	if req.FBv > 0 {
		msg.FBv = req.FBv
	} else {
		vk := fmt.Sprintf("BV_%d", head.Uid)
		bv, _ := h.userRedisMgr.RedisGet(vk)
		msg.FBv, _ = strconv.Atoi(bv)
	}

	msgBuf := h.SendMsg2Mq(head, msg, tail)

	if err := msg.SaveUserMsg(msgBuf); err != nil {
		logger.Error(err)
	}

	respStr := fmt.Sprintf(`{"msgid":%d,"sendtime":%d}`, msg.MsgId, msg.SendTime)

	return []byte(respStr), 0
}

func (h *Handler) ReceivedMsg(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req common.MsgRecvedUser

	if err := json.Unmarshal(jsonBody, &req); err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	logger.Info("remove msg touid:", head.Uid, " msgid:", req.MsgId)

	result := 0

	if common.IsChatMsgId(req.MsgId) {
		userMsg := &common.UserMsgItem{0, req.MsgId, req.Uid, head.Uid, 0, "", 0, "", 0, 0, ""}

		if err := userMsg.DelUserMsg(); err != nil {
			logger.Error(err)
		}
	} else {
		msgKey := fmt.Sprintf("%s%d", common.SET_TEAMMSGID, head.Uid)
		val := fmt.Sprintf("%d", req.MsgId)
		h.teamMsgRedisMgr.RedisZRem(msgKey, val)
	}

	if result != 0 {
		return []byte(""), common.ERR_CODE_SYS
	} else {
		return []byte(""), 0
	}
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	h.count++
	msgBody := bytes.NewReader(message.Body)
	head := common.PkgHead{}
	err := binary.Read(msgBody, binary.BigEndian, &head)

	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return err
	}

	if int(head.PkgLen)+common.SIZEOF_INNERTAIL != len(message.Body) {
		logger.Error("pkghead len(", head.PkgLen, ") , pkg len(", len(message.Body), ") is error")
		return nil
	}

	jsonStr := make([]byte, head.PkgLen-common.SIZEOF_PKGHEAD)
	err = binary.Read(msgBody, binary.BigEndian, &jsonStr)
	if err != nil {
		logger.Error("read req json fail:", err.Error())
		return err
	}

	logger.Info("req json:", string(jsonStr[:]))

	tail := common.InnerPkgTail{}
	err = binary.Read(msgBody, binary.BigEndian, &tail)

	if err != nil {
		logger.Error("read req innertail fail:", err.Error())
		return err
	}

	logger.Debug("req head:", head.ToString())
	logger.Debug("req tail:", tail.ToString())

	var resp []byte
	var retCode uint32 = 0

	switch head.Cmd {
	case common.DU_CMD_IM_SEND_USER_MSG:
		resp, retCode = h.SendMsg(head, jsonStr, &tail)
	case common.DU_CMD_IM_USER_MSG_RECEIVED:
		resp, retCode = h.ReceivedMsg(head, jsonStr, &tail)
	default:
		logger.Debug(" error Cmd :", head.Cmd)
		return nil
	}

	logger.Info("msg_server HandleMessage() Cmd", head.Cmd, " resp:", string(resp[:]), " retcode:", retCode)

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
		err := h.producers[h.count%(len(h.producers))].Publish(topic, respBuf.Bytes())
		if err != nil {
			logger.Error("----------Publish", topic, "err:", err.Error())
		} else {
			logger.Debug("---------------Publish ", topic, " success")
		}
	}

	return nil
}

func SendMsg2MsgCenter(msgBodyCh chan []byte) int {
	producer, err := nsq.NewProducer(*g_nsqd_addrs, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	defer producer.Stop()

	for {
		select {
		// Publish “db2msg_topic”
		case msg := <-msgBodyCh:
			err = producer.Publish(*g_Db2MsgTopic, msg)
			if err != nil {
				logger.Error("--------- Publish ", *g_Db2MsgTopic, "err:", err.Error())
			} else {
				logger.Debug("--------- Publish", *g_Db2MsgTopic, " success")
			}
		}
	}

	return 0
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

	teamRedis := common.NewRedisManager(*g_RedisAddr)
	if teamRedis == nil {
		logger.Info("connect team redis", *g_RedisAddr, "fail")
		return
	}
	logger.Info("connect  team redis ", *g_RedisAddr, " success")

	userRedis := common.NewRedisManager(*g_UserRedisAddr)
	if userRedis == nil {
		logger.Info("connect user redis ", *g_UserRedisAddr, "fail")
		return
	}
	logger.Info("connect  user redis ", *g_UserRedisAddr, " success")

	msgRedis := common.NewRedisManager(*g_TeamMsgRedisAddr)
	if msgRedis == nil {
		logger.Info("connect user redis ", *g_TeamMsgRedisAddr, "fail")
		return
	}
	logger.Info("connect  msg redis ", *g_TeamMsgRedisAddr, " success")

	//	kefuMsgRedis := common.NewRedisManager(*g_KefuMsgRedisAddr)
	//	if msgRedis == nil {
	//		logger.Info("connect user redis ", *g_KefuMsgRedisAddr, "fail")
	//		return
	//	}
	//	logger.Info("connect kefu msg redis ", *g_KefuMsgRedisAddr, " success")

	if common.MongoInit(*g_MongodbAddr) != 0 {
		logger.Info("connect mongo ", *g_MongodbAddr, "fail")
		return
	}
	logger.Info("connect mongo ", *g_MongodbAddr, " success")

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

	r, err := nsq.NewConsumer(*g_Conn2MsgSvrTopic, "msg_server-channel", nsq.NewConfig())
	if err != nil {
		logger.Error("NewConsumer ", *g_Conn2MsgSvrTopic, " error:", err.Error())
		return
	}

	defer r.Stop()

	r.SetLogger(nil, nsq.LogLevelInfo)

	ch := make(chan int)

	chatMsgCh := make(chan []byte, 1000)

	go SendMsg2MsgCenter(chatMsgCh)

	handler := &Handler{0, ch, chatMsgCh, producers, teamRedis, userRedis, msgRedis}

	r.AddHandler(handler)

	addrs = strings.Split(*g_nsqd_addrs, ",")
	err = r.ConnectToNSQDs(addrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	fmt.Println("StartServer")
	logger.Info("StartServer")

	select {}

	return
}
