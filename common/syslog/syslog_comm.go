package syslog

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
	sync_ "sirendaou.com/duserver/common/sync"
)

const (
	ALL_LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type LogMsg struct {
	Name     string
	Time     time.Time
	Filename string
	Level    int
	Content  string
}

type logMgr struct {
	waitGroup     *sync_.WaitGroup
	logTopic      string
	producers     []*nsq.Producer
	producerCount uint
	logMsg        chan []byte
	logLevel      int
	modelName     string
}

type Config struct {
	NsqdAddrs string
	LogTopic  string
	ModelName string
	LogLevel  int
}

var (
	g_sysLogger     *logMgr = nil
	g_defaultConfig         = &Config{
		NsqdAddrs: "127.0.0.1:4150",
		LogTopic:  "sysLogTopic",
		ModelName: "UnKnow",
	}
)

func SysLogInit(config *Config) error {
	if g_sysLogger != nil {
		log.Println("syslog already init once! cat't init again")
		return nil
	}
	var nsqdAddrs, logTopic, modelName string
	var logLevel int = 0
	if config != nil && len(config.NsqdAddrs) > 0 {
		nsqdAddrs = config.NsqdAddrs
	} else {
		nsqdAddrs = g_defaultConfig.NsqdAddrs
	}

	if config != nil && len(config.LogTopic) > 0 {
		logTopic = config.LogTopic
	} else {
		logTopic = g_defaultConfig.LogTopic
	}

	if config != nil && len(config.ModelName) > 0 {
		modelName = config.ModelName
	} else {
		modelName = g_defaultConfig.ModelName
	}

	if config != nil {
		logLevel = config.LogLevel
	}

	addrSlice := strings.Split(nsqdAddrs, ",")
	producerCount := len(addrSlice)
	producers := make([]*nsq.Producer, producerCount)

	var err error
	for i, addr := range addrSlice {
		config := nsq.NewConfig()
		config.DefaultRequeueDelay = 0
		producers[i], err = nsq.NewProducer(addr, config)

		if err != nil {
			fmt.Println("NewProducer ", addr, " error:", err)
			return err
		}
	}

	if len(logTopic) == 0 {
		logTopic = "sysLogTopic"
	}
	g_sysLogger = &logMgr{
		waitGroup:     sync_.NewWaitGroup(),
		producers:     producers,
		producerCount: uint(producerCount),
		logMsg:        make(chan []byte),
		logTopic:      logTopic,
		modelName:     modelName,
		logLevel:      logLevel,
	}

	go g_sysLogger.run()
	return nil
}

func (l *logMgr) run() {
	l.waitGroup.AddOne()
	defer l.waitGroup.Done()

	count := uint(0)
	for {
		select {
		case <-l.waitGroup.ExitNotify():
			return
		case msg := <-l.logMsg:
			count++
			err := l.producers[count%l.producerCount].Publish(l.logTopic, msg)
			fmt.Println("Publish Msg:", string(msg))
			if err != nil {
				fmt.Println("nsq.Publish error:", err)
			}
		}
	}
}

func sysLog(logMsg *LogMsg) {
	body, err := json.Marshal(logMsg)
	if err != nil {
		log.Println("Json.Marshal failed:", err)
		return
	}
	g_sysLogger.logMsg <- body
}

func SysLogDeinit() {
	g_sysLogger.waitGroup.Wait()

	for _, produce := range g_sysLogger.producers {
		produce.Stop()
	}
}

func NewLogMsg(level int, content string) *LogMsg {
	_, file, line, _ := runtime.Caller(2)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	logMsg := &LogMsg{
		Name:     g_sysLogger.modelName,
		Time:     time.Now(),
		Filename: fmt.Sprint(short, ":", line),
		Level:    level,
		Content:  content,
	}
	return logMsg
}

func levelString(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	case ALL_LEVEL:
		return "ALL_LEVEL"
	case OFF:
		return "OFF"
	}
	return "UNKNOW"
}

func (logMsg *LogMsg) Format() string {
	return fmt.Sprintln(logMsg.Time.Format("2006/01/02 15:04:05"), " ", logMsg.Filename, " ", logMsg.Name, " ", levelString(logMsg.Level), " ", logMsg.Content)
}

func console(logMsg *LogMsg) {
	fmt.Println(logMsg.Format())
}

func Debug(v ...interface{}) {
	logMsg := NewLogMsg(DEBUG, fmt.Sprint(v))
	if g_sysLogger.logLevel <= DEBUG {
		console(logMsg)
		sysLog(logMsg)
	}
}

func Info(v ...interface{}) {
	logMsg := NewLogMsg(INFO, fmt.Sprint(v))
	if g_sysLogger.logLevel <= INFO {
		console(logMsg)
		sysLog(logMsg)
	}
}

func Warn(v ...interface{}) {
	logMsg := NewLogMsg(WARN, fmt.Sprint(v))
	if g_sysLogger.logLevel <= WARN {
		console(logMsg)
		sysLog(logMsg)
	}
}

func Error(v ...interface{}) {
	logMsg := NewLogMsg(ERROR, fmt.Sprint(v))
	if g_sysLogger.logLevel <= ERROR {
		console(logMsg)
		sysLog(logMsg)
	}
}

func Fatal(v ...interface{}) {
	logMsg := NewLogMsg(FATAL, fmt.Sprint(v))
	if g_sysLogger.logLevel <= FATAL {
		console(logMsg)
		sysLog(logMsg)
	}
}
