package syslog_server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/bitly/go-nsq"
	"github.com/rakyll/globalconf"

	"sirendaou.com/duserver/common/syslog"
)

var (
	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
	g_log_file  = flag.String("log_file", "syslog_server.log", "the log file path")
	g_log_level = flag.Int("log_level", 2, "the level of log")

	g_cpu_num      = flag.Int("cpu_num", 4, "the num of cpu")
	g_nsqdTCPAddrs = flag.String("nsq_addr", "", "nsq Server address (transient)")
	g_sysLogTopic  = flag.String("syslog_topic", "sysLogTopic", "the name of syslog to msg center")
)

type Handler struct {
	log *logFile
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	logMsg := &syslog.LogMsg{}
	if err := json.Unmarshal(message.Body, logMsg); err != nil {
		log.Println("json.Unmarshal Failed:", err)
		return err
	}
	h.log.WriteLogMsg(logMsg)
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

	// init logger
	logConfig := &syslog.Config{
		NsqdAddrs: *g_nsqdTCPAddrs,
		LogTopic:  *g_sysLogTopic,
		ModelName: "syslog_server",
	}
	syslog.SysLogInit(logConfig)

	if len(*g_nsqdTCPAddrs) < 1 {
		flag.PrintDefaults()
		return
	}

	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0

	consumer, err := nsq.NewConsumer(*g_sysLogTopic, "syslog-channel", nsq.NewConfig())
	if err != nil {
		syslog.Error("NewConsumer ", *g_sysLogTopic, " error:", err.Error())
		return
	}
	syslog.Info("NewConsumer ", *g_sysLogTopic, " ok:")

	defer consumer.Stop()

	consumer.SetLogger(nil, nsq.LogLevelInfo)

	// 新建日志处理器
	log := NewLogFile(*g_log_path, *g_log_file)
	handler := &Handler{log: log}
	consumer.AddHandler(handler)

	addrs := strings.Split(*g_nsqdTCPAddrs, ",")
	err = consumer.ConnectToNSQDs(addrs)
	if err != nil {
		fmt.Println("ConnectToNSQ err:", err.Error())
		return
	}

	select {}

	return
}
