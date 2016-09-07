package sis

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/donnie4w/go-logger/logger"
	"github.com/rakyll/globalconf"
)

const (
	UDP_PORT = 19000
)

var (
	g_ConnAddr = flag.String("connaddr", "", "conn server addr")
	g_udpPort  = flag.Int("udpport", 0, "sis udp port addr")

	g_log_path  = flag.String("log_path", "./logs", "the log file path directory")
	g_log_file  = flag.String("log_file", "sis_server.log", "the log file name")
	g_log_level = flag.Int("log_level", 2, "the log level 1-debug 2-info(default) 3-WARN 4-error 5-FATAL 6-off")
)

type SisInfo struct {
	Len    int16
	Head   [2]byte
	Net    [30]byte
	Opear  int32
	Uid    uint32
	Appkey [50]byte
	Ver    [10]byte
	Mode   int32
	Res    [22]byte
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

//	logger.SetConsole(false)
	logger.SetRollingDaily(*g_log_path, *g_log_file)
	logger.SetLevel(logger.LEVEL(*g_log_level))

	logger.Info("StartServer g_log_path=", *g_log_path)
	logger.Info("StartServer g_log_file=", *g_log_file)
	logger.Info("StartServer g_log_level=", *g_log_level)

	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: *g_udpPort,
	})

	if err != nil {
		logger.Error("bind fail", err)
		return
	}

	defer socket.Close()

	respjspntest := `{"ips":["192.168.20.179:9100"], "file":"192.168.20.179:8889"}`
	respjspn := fmt.Sprintf(`{"ips":["%s"], "file":"192.168.20.179:9100"}`, *g_ConnAddr)

	req := SisInfo{}
	for {
		data := make([]byte, 256)
		_, toaddr, err := socket.ReadFromUDP(data)

		if err != nil {
			logger.Error("read fail", err)
			continue
		}

		logger.Info("toaddr=", toaddr)

		p := bytes.NewReader(data)

		err = binary.Read(p, binary.BigEndian, &req)

		if err != nil {
			logger.Error("Read SisInfo err:", err)
		} else {
			logger.Info("uid:", req.Uid, "appkey:", string(req.Appkey[0:]), "Net:", string(req.Net[0:]), "Opr:", req.Opear, "Mode:", req.Mode)
		}

		if req.Mode == 0 {
			socket.WriteTo([]byte(respjspn), toaddr)
		} else {
			socket.WriteTo([]byte(respjspntest), toaddr)
		}
	}

}
