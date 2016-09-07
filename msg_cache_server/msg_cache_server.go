package msg_cache_server

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"

	"sirendaou.com/duserver/common"
	"errors"
)

var (
	g_log_path = flag.String("log_path", "/tmp", "the log file path")
	g_log_file = flag.String("log_file", "loadreg2es.log", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")
	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")

	g_RedisAddr = flag.String("redis_addr", "", "redis mq server addr")
	g_nsqdTCPAddrs = flag.String("nsq_addr", "", "nsq Server address (transient)")

	g_Db2MsgTopic = flag.String("db2msg_topic", "t-db2msgcenter", "the name of db to msg center")
)

type Handler struct {
	count int
	producers []*nsq.Producer
	redisMgr *common.RedisManager
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	h.count++

	ret, head, jsonStr, tail := common.DecPkgInnerBody(message.Body)
	if ret != 0 {
		logger.Error("DecPkgInnerBody fail ,ret ", ret)
		return errors.New("DecPkgInnerBody error")
	}

	logger.Debug("head:",head.ToString())
	logger.Debug("tail:",tail.ToString())

//	common.MsgDBFlagSet(tail.ToUid)
	if err := h.redisMgr.RedisMsgCacheSet(tail.ToUid, "unread"); err != nil {
		logger.Error(err)
		return err
	}

	switch head.Cmd {
	case common.DU_CMD_IM_SEND_USER_MSG:
		head.Cmd = common.DU_PUSH_CMD_IM_USER_MSG
	case common.DU_PUSH_CMD_IM_SYSTEM_MSG:
		logger.Info("Push sys msg")
		break
	case common.DU_PUSH_CMD_IM_USER_MSG:
		logger.Info("Push im msg")
		break
	case common.DU_PUSH_CMD_IM_TEAM_MSG:
		logger.Info("Push team msg")
		break
	default:
		logger.Debug("error Cmd :", head.Cmd)
	}

	toTail, err := h.redisMgr.RedisStatCacheGet(tail.ToUid)
	if err != nil || toTail.ConnIP == 0 || toTail.ConnPort == 0 || toTail.Sid == 0 || toTail.ToUid == 0 {
		logger.Debug("------------------------not online!!!!!!!!")
		return err
	} else {
		logger.Debug("StatCacheGet:", toTail.ToString())
	}

	toTail.FromUid = head.Uid
	toTail.MsgId = tail.MsgId
	toTail.ToUid = tail.ToUid

	logger.Debug("toTail=", toTail.ToString())

	respBuf := new(bytes.Buffer)

	binary.Write(respBuf, binary.BigEndian, head)
	binary.Write(respBuf, binary.BigEndian, jsonStr)
	binary.Write(respBuf, binary.BigEndian, toTail)

	topic := fmt.Sprintf("conn_%s_%d", common.IntToIP(toTail.ConnIP), toTail.ConnPort)
	//resp send to conn
	if len(topic) > 0 {
		err := h.producers[h.count%len(h.producers)].Publish(topic, respBuf.Bytes())
		if err != nil {
			logger.Error("nsq write resp to ", topic, "err:", err.Error())
			return err
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

	if len(*g_nsqdTCPAddrs) < 1 {
		flag.PrintDefaults()
		return
	}

	redisMgr := common.NewRedisManager(*g_RedisAddr)
	if redisMgr == nil {
		logger.Info("connect redis ", *g_RedisAddr, "fail")
		return
	}

	addrs := strings.Split(*g_nsqdTCPAddrs, ",")
	producerCount := len(addrs)
	producers := make([]*nsq.Producer, producerCount)
	for i, addr := range addrs {
		producers[i], err = nsq.NewProducer(addr, nsq.NewConfig())
		if err != nil {
			logger.Error("NewProducer  ", addr, " err:", err)
			return
		}
	}

	defer func() {
		for _, p := range producers {
			p.Stop()
		}
	}()

	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0

	r, err := nsq.NewConsumer(*g_Db2MsgTopic, "msgcenter-channel", nsq.NewConfig())
	if err != nil {
		logger.Error("NewConsumer ", *g_Db2MsgTopic, " error:", err)
		return
	}
	logger.Info("NewConsumer ", *g_Db2MsgTopic, "ok")

	defer r.Stop()

	r.SetLogger(nil, nsq.LogLevelInfo)

	ch := make(chan int)

	handler := &Handler{0, producers,redisMgr}
	r.AddHandler(handler)

	addrs = strings.Split(*g_nsqdTCPAddrs, ",")

	err = r.ConnectToNSQDs(addrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err)
		return
	}

	select {
	case ret := <-ch:
		fmt.Println("quit:", ret)
	}

	return
}
