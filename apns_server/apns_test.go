package apns_server

import (
	"fmt"
	"os"
	"sirendaou.com/duserver/apns_server/apns"
	"github.com/rakyll/globalconf"
	"flag"
)

var (
//	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
//	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
//	g_log_path  = flag.String("log_path", "/tmp", "the log file path")
//	g_log_path  = flag.String("log_path", "/tmp", "the log file path")

	g_PemFilePath = flag.String("pemfile_path", "", "the path of pem  file")
	g_KeyFilePath = flag.String("keyfile_path", "", "the path of pem  file")
)

func main() {
//	if len(os.Args) < 5 {
//		fmt.Println(os.Args[0], "pemfile, keyfile , token , mod(1-sandbox 2-product")
//		return
//	}
//
//	keys := os.Args[1]
//	pem := os.Args[2]

	conf, err := globalconf.NewWithOptions(&globalconf.Options{
		Filename: os.Args[1],
	})
	if err != nil {
		fmt.Print("NewWithFilename ", os.Args[1], " fail :", err)
		return
	}

	conf.ParseAll()

	payload := apns.NewPayload()
	payload.Alert = "test!"
	payload.Badge = 42
	payload.Sound = "bingbong.aiff"

	pn := apns.NewPushNotification()
//	pn.DeviceToken = os.Args[3]
	pn.DeviceToken = ""
	pn.AddPayload(payload)

	mod := os.Args[4]

	url := ""
	switch mod {
	case "1":
		url = "gateway.sandbox.push.apple.com:2195"
	default:
		url = "gateway.push.apple.com:2195"
	}

	client := apns.NewClient(url, *g_PemFilePath, *g_KeyFilePath)

	resp := client.Send(pn)

	alert, _ := pn.PayloadString()

	fmt.Println("  Alert:", alert)
	fmt.Println("Success:", resp.Success)
	fmt.Println(*g_PemFilePath, url, "  Error:", resp.Error)
}
