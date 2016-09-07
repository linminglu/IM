package file_server

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"
	"github.com/rakyll/globalconf"
)

var (
	g_ListenAddr = flag.String("listen_addr", "0.0.0.0:8889", "http server addr")

	g_log_path        = flag.String("log_path", "/tmp", "the log file path")
	g_log_file        = flag.String("log_file", "file_server.log", "the log file path")
	g_log_level       = flag.Int("log_level", 2, "the log level 1-debug 2-info(default) 3-WARN 4-error 5-FATAL 6-off")

	g_cpu_num = flag.Int("cpu_num", 4, "the num of cpu")

	g_CallbackUrl = flag.String("callback_url", "", "qiniu upload callback url")
	g_DownloadUrl = flag.String("download_url", "", "qiniu download url")
	g_UploadUrl   = flag.String("upload_url", "", "qiniu upload url")
)

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

	logger.Info("file server StartServer")

//	logfile, err := os.OpenFile(*g_file_server_log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		fmt.Printf("%s\r\n", err.Error())
//		os.Exit(-1)
//	}

//	defer logfile.Close()
//
//	svrLog := log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

	web.Post("(/file/upload)", UploadHandler)
	web.Get("(/file/upload)", uploadHandlerGet)
	web.Get("(/file/upload2)", uploadHandlerGet2)
	web.Post("(/file/callback)", CallBack)
	web.Get("(/file/download)", downloadHandler)
	web.Post("(/file/download)", downloadHandlerPost)
	web.Post("(/file/image)", imageHandler)

//	web.SetLogger(svrLog)

	web.Run(*g_ListenAddr)
}
