package apns_server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"

	"sirendaou.com/duserver/common"
)

var (
	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
	g_log_name  = flag.String("log_file", "loadreg2es.log", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_cpunum = flag.Int("cpu_num", 4, "the num of cpu")
	//	g_isdaemon      = flag.Int("isdaemon", 0, "is run server as daemon 0-no 1-s)")
	//	lifetime        = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
	//	g_shmkey        = flag.Int("shmkey", 22222, "stat center user stat info shm key)")

	g_AppRedisAddr  = flag.String("app_redis_addr", "", "appinfo redis server addr")
	g_UserRedisAddr = flag.String("user_redis_addr", "", "userinfo redis mq server addr")

	g_PemFilePath = flag.String("pemfile_path", "", "the path of pem  file")

	g_mysql_host = flag.String("mysql_host", "", "mysql host")
	g_mysql_db   = flag.String("mysql_db", "", "mysql db name")
	g_mysql_user = flag.String("mysql_user", "", "mysql user")
	g_mysql_pwd  = flag.String("mysql_pwd", "", "mysql passwd")

	g_nsqdTCPAddrs = flag.String("nsq_addr", "", "nsq Server address (transient)")
	g_2ApnsTopic   = flag.String("msg2apns_topic", "t-msg2apns", "the name of msgcenter to apns topic")
)

var g_ApnsProcMap map[string]*ApnsProc

type UserMsgReq struct {
	MsgContent string `json:"msgcontent,omitempty"`
	ToUid      uint64 `json:"touid,omitempty"`
	MsgType    uint16 `json:"msgtype,omitempty"`
	ApnsText   string `json:"apnstext,omitempty"`
}

type Handler struct {
	reqCh chan common.ApnsMsg
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	reader := bytes.NewReader(message.Body)
	head := common.PkgHead{}

	err := binary.Read(reader, binary.BigEndian, &head)
	if err != nil {
		logger.Error("read pkghead fail:", err.Error())
		return nil
	}

	logger.Debug("PkgHead:", head.ToString())

	if head.Cmd != common.DU_PUSH_CMD_IM_USER_MSG && common.DU_PUSH_CMD_IM_TEAM_MSG != head.Cmd && common.DU_PUSH_CMD_IM_SYSTEM_MSG != head.Cmd {
		logger.Debug("not support apns cmd=", head.Cmd)
		return nil
	}

	if int(head.PkgLen)+common.SIZEOF_INNERTAIL != len(message.Body) {
		logger.Error("pkghead len(", head.PkgLen, ") , pkg len(", len(message.Body), ") is error")
		return nil
	}

	jsonStr := make([]byte, head.PkgLen-common.SIZEOF_PKGHEAD)

	err = binary.Read(reader, binary.BigEndian, &jsonStr)
	if err != nil {
		logger.Error("read req json fail:", err.Error())
		return nil
	}

	logger.Info("req json:", string(jsonStr[:]))

	tail := common.InnerPkgTail{}
	err = binary.Read(reader, binary.BigEndian, &tail)
	if err != nil {
		logger.Error("read req innertail fail:", err.Error())
		return nil
	}

	logger.Info("req innertail:", tail.ToString())
	logger.Info("req json:", string(jsonStr[:]))

	apnsReq := common.ApnsMsg{}
	msgReq := UserMsgReq{}

	switch head.Cmd {
	case common.DU_PUSH_CMD_IM_USER_MSG:
	case common.DU_PUSH_CMD_IM_TEAM_MSG:
		err := json.Unmarshal(jsonStr, &msgReq)
		if err != nil {
			logger.Error("Unmarshal error:", err)
			return nil
		}

		if head.Cmd != common.DU_PUSH_CMD_IM_USER_MSG && msgReq.ToUid == 0 {
			msgReq.ToUid = tail.FromUid
		}

		if len(msgReq.ApnsText) < 1 {
			logger.Info("msgReq.Apnstext is empty")
			return nil
		}

		if msgReq.ToUid&0x0f != common.PT_IOS {
			logger.Error("uid ", msgReq.ToUid, " is not ios")
			return nil
		}

		apnsReq.Msg = msgReq.ApnsText
		//		apnsReq.AppKey = ""
		apnsReq.UidList = make([]uint64, 1, 1)
		apnsReq.UidList[0] = msgReq.ToUid

		// -> reqCh -> ProcApnsReq()
		h.reqCh <- apnsReq
	default:
		logger.Debug(" error Cmd :", head.Cmd)
		return nil
	}

	return nil
}

type AppStageTime struct {
	//	AppKey    string
	LTime  int64
	NStage int
}

var g_AppStageMap map[string]AppStageTime = nil

func getStage(appRedis *common.RedisManager) int {
	if g_AppStageMap == nil {
		g_AppStageMap = make(map[string]AppStageTime)
	}

	keys := "appinfo_du"

	st, ok := g_AppStageMap[keys]

	now := time.Now().Unix()
	if ok {
		if st.LTime+120 >= now {
			return st.NStage
		}
	}

	stage, errcode := appRedis.RedisHGet(keys, "ios_stage")
	iosStage := 1 // test is default

	logger.Info("RedisHGet ", keys, " ios_stage", stage, errcode)

	var err error
	if errcode == 0 {
		iosStage, err = strconv.Atoi(stage)
		if err != nil {
			iosStage = 1
		} else {
			g_AppStageMap[keys] = AppStageTime{now, iosStage}
		}
	}

	return iosStage
}

func ProcApnsReq(reqCh chan common.ApnsMsg, userRedis, appRedis *common.RedisManager, quitCh chan string) {
	for {
		select {
		case quitAppKey := <-quitCh:
			logger.Info(quitAppKey, "quit")
			delete(g_ApnsProcMap, quitAppKey)
		case apnsReq := <-reqCh:
			if len(apnsReq.UidList) < 1 {
				logger.Error("len( apnsreq.UidList) == 0 ")
				continue
			}

			logger.Info("apnsReq:", apnsReq)

			//			if len(apnsReq.AppKey) != 24 {
			//				// get appkey by uid
			//				appId, _ := common.GetAppIdPtFromUid(apnsReq.UidList[0])
			//				appKey, ok := g_AppKeyMap[int(appId)]
			//				if !ok {
			//					keys := fmt.Sprintf("appid_%d", appId)
			//					appKey, errcode := AppRedis.RedisGet(keys)
			//					if errcode == 0 && len(appKey) == 24 {
			//						g_AppKeyMap[int(appId)] = appKey
			//						apnsReq.AppKey = appKey
			//					} else {
			//						logger.Error("Uid:", apnsReq.UidList[0], "appid:", appId, "cannt find appkey")
			//					}
			//				} else {
			//					apnsReq.AppKey = appKey
			//				}
			//			}
			//
			//			logger.Info("apnsreq:", apnsReq)
			//
			//			if len(apnsReq.AppKey) != 24 {
			//				logger.Error("cannot get right appkey:", apnsReq.AppKey)
			//				continue
			//			}

			stage := getStage(appRedis)
			key := fmt.Sprintf("%s_%d", "du", stage)
			proc, ok := g_ApnsProcMap[key]
			if ok {
				proc.sendMsg(apnsReq)
			} else {
				// add new proc
				apnsProc := NewApnsProc(*g_PemFilePath, stage, quitCh, userRedis, key)
				g_ApnsProcMap[key] = apnsProc

				logger.Info("add apns proc :", key)

				// start proc
				go apnsProc.StartProc()

				apnsProc.sendMsg(apnsReq)
			}
		}
	}

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

	runtime.GOMAXPROCS(*g_cpunum)

	//	logger.SetConsole(false)
	logger.SetRollingDaily(*g_log_path, *g_log_name)
	logger.SetLevel(logger.LEVEL(*g_log_level))

	if common.MysqlInit(*g_mysql_host, *g_mysql_db, *g_mysql_user, *g_mysql_pwd) != 0 {
		logger.Error("mysql init ", *g_mysql_host, *g_mysql_db, *g_mysql_user, "fail")
		return
	}

	if len(*g_nsqdTCPAddrs) < 1 {
		flag.PrintDefaults()
		return
	}

	g_ApnsProcMap = make(map[string]*ApnsProc, 1000)

	userRedis := common.NewRedisManager(*g_UserRedisAddr)
	if userRedis == nil {
		logger.Error("connect user redis ", *g_UserRedisAddr, "fail")
		return
	}

	appRedis := common.NewRedisManager(*g_AppRedisAddr)
	if appRedis == nil {
		logger.Error("connect app redis ", *g_AppRedisAddr, "fail")
		return
	}

	go FeedBackProc(userRedis)

	// apnsReqCh
	apnsReqCh := make(chan common.ApnsMsg, 1000)
	procQuitCh := make(chan string, 1)

	go ProcApnsReq(apnsReqCh, userRedis, appRedis, procQuitCh)

	consumer, err := nsq.NewConsumer(*g_2ApnsTopic, "apnssvr-channel", nsq.NewConfig())
	if err != nil {
		logger.Error("NewConsumer ", *g_2ApnsTopic, " error:", err.Error())
		return
	}
	defer consumer.Stop()
	consumer.SetLogger(nil, nsq.LogLevelInfo)

	// consumer.AddHandle
	handler := &Handler{reqCh: apnsReqCh}
	consumer.AddHandler(handler)

	err = consumer.ConnectToNSQD(*g_nsqdTCPAddrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	procStopCh := make(chan string)
	for {
		select {
		case appKey := <-procStopCh:
			delete(g_ApnsProcMap, appKey)
		}
	}

	return
}
