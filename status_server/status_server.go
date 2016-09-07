package status_center

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

	"sirendaou.com/duserver/common"
)

var (
	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
	g_log_file  = flag.String("log_file", "loadreg2es.log", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_cpu_num  = flag.Int("cpu_num", 4, "the num of cpu")

	g_RedisAddr = flag.String("redis_addr", "", "redis mq server addr")
	g_nsqdTCPAddrs   = flag.String("nsq_addr", "", "nsq Server address (transient)")
	g_Conn2StatTopic = flag.String("conn2stat_topic", "t-conn2statcenter", "the name of connect to status center")
)

type Handler struct {
	count    int
	producer []*nsq.Producer
	redisMgr *common.RedisManager
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	h.count++

	reader := bytes.NewReader(message.Body)

	head := common.InnerPkgTail{}

	err := binary.Read(reader, binary.BigEndian, &head)
	if err != nil {
		logger.Error(err)
		return err
	}

	// is login Conflict
//	if head.ToUid == 1 {
//		tail, err := h.redisMgr.RedisStatCacheGet(head.FromUid)
//		if err != nil {
//			logger.Error(err)
//			return err
//		}
//
//		if tail.ToUid > 0 && tail.Sid > 0 && tail.ConnIP > 0 || tail.ConnPort > 0 {
//			respBuf := new(bytes.Buffer)
//
//			loginConflictHead := common.PkgHead{common.SIZEOF_PKGHEAD, common.DU_PUSH_CMD_USER_LOGIN_CONFLICT, 1, 0, tail.Sid, tail.FromUid, 0}
//			loginConflictTail := common.InnerPkgTail{tail.ConnIP, tail.ConnPort, head.FromUid, 0, tail.Sid, 0}
//
//			binary.Write(respBuf, binary.BigEndian, loginConflictHead)
//			binary.Write(respBuf, binary.BigEndian, loginConflictTail)
//
//			topic := fmt.Sprintf("conn_%s_%d", common.IntToIP(tail.ConnIP), tail.ConnPort)
//			//resp send to conn
//			err := h.producer[h.count%len(h.producer)].Publish(topic, respBuf.Bytes())
//			if err != nil {
//				logger.Error("nsq write loginConflict to ", topic, "err:", err.Error())
//			} else {
//				logger.Debug("---------------Publish ", topic, " success")
//			}
//		}
//	}

	head.ToUid = head.FromUid

	if err = h.redisMgr.RedisStatCacheSet(head); err != nil {
		logger.Error(err)
		return err
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

	config := nsq.NewConfig()
	// so that the test can simulate reaching max requeues and a call to LogFailedMessage
	config.DefaultRequeueDelay = 0
	// so that the test wont timeout from backing off
	//config.MaxBackoffDuration = time.Millisecond * 50
	//config.Deflate = false
	//config.Snappy =  false

	//	common.StatCacheInit(uint32(*g_shmkey))

	consumer, err := nsq.NewConsumer(*g_Conn2StatTopic, "statcenter-channel", nsq.NewConfig())
	if err != nil {
		logger.Error("NewConsumer ", *g_Conn2StatTopic, " error:", err.Error())
		return
	}

	defer consumer.Stop()

	redisMgr := common.NewRedisManager(*g_RedisAddr)
	if redisMgr == nil {
		logger.Info("connect redis ", *g_RedisAddr, "fail")
		return
	}

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

	consumer.SetLogger(nil, nsq.LogLevelInfo)

	handler := &Handler{0, producers, redisMgr}
	consumer.AddHandler(handler)

	addrs = strings.Split(*g_nsqdTCPAddrs, ",")
	err = consumer.ConnectToNSQDs(addrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	select {}

	return
}
