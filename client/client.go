package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"IM/common/json"
	sync_ "IM/common/sync"

	"github.com/astaxie/beego/logs"
)

var (
	g_client = &Client{
		ReqMsgCh:  make(chan *json.Json, 10),
		RespMsgCh: make(chan *json.Json, 10),
		waitGroup: sync_.NewWaitGroup(),
	}
)

type Client struct {
	Conn net.Conn

	ReqMsgCh  chan *json.Json
	RespMsgCh chan *json.Json
	waitGroup *sync_.WaitGroup
}

func (this *Client) readNetData() (*json.Json, int) {
	var recvBuf [MAX_RECVBUF_SIZE]byte
	this.Conn.SetDeadline(time.Now().Add(NET_READ_DEADLINE * time.Second))
	if _, err := io.ReadFull(this.Conn, recvBuf[:2]); err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			return nil, IM_ERR_NET_TIMEOUT
		} else if !ok {
			logs.Debug("read net data error:", err)
			return nil, IM_ERR_NETWORK
		}
	}

	dataSize := binary.BigEndian.Uint16(recvBuf[:2])

	if dataSize > MAX_RECVBUF_SIZE || dataSize < 0 {
		logs.Debug("data package size error:", dataSize)
		return nil, IM_ERR_DATA_PACKAGE_SIZE
	}

	this.Conn.SetDeadline(time.Now().Add(NET_READ_DEADLINE * time.Second))
	if _, err := io.ReadFull(this.Conn, recvBuf[:dataSize]); err != nil {
		logs.Debug("read net data failed:", err)
		return nil, IM_ERR_NETWORK
	}

	jsonObj, err := json.Unmarshal(recvBuf[:dataSize])
	if err != nil {
		logs.Debug("json parse failed:", err)
		return nil, IM_ERR_DATA_FORMAT
	}

	return jsonObj, IM_ERR_SUCCESS
}

func (this *Client) writeNetData(data *json.Json) {
	buf := new(bytes.Buffer)

	body, err := json.Marshal(data)
	if err != nil {
		logs.Debug(err, data.ToString())
		return
	}

	dataSize := len(body)
	binary.Write(buf, binary.BigEndian, dataSize)
	binary.Write(buf, binary.BigEndian, body)

	this.Conn.SetDeadline(time.Now().Add(NET_WRITE_DEADLINE * time.Second))
	_, err = this.Conn.Write(buf.Bytes())

	if err != nil {
		logs.Debug(err)
	}
}

func (this *Client) Serve() {
	go this.Handle()

	lastTime := time.Now().Unix()
	err := IM_ERR_NET_TIMEOUT
	for lastTime+NET_DISCONNECT_DEADLINE > time.Now().Unix() {
		data, err := this.readNetData()

		if err == IM_ERR_NET_TIMEOUT {
			continue
		} else if err == IM_ERR_SUCCESS {
			lastTime = time.Now().Unix()
			this.ReqMsgCh <- data
		} else {
			break
		}
	}

	logs.Info("client exit code:", err, " addr:", this.Conn.RemoteAddr())

	this.waitGroup.Wait() // 等待处理协程结束
	this.Conn.Close()     // 关闭连接

}

func (this *Client) Handle() {
	this.waitGroup.Add(1)
	defer this.waitGroup.Done()

	exitNotify := this.waitGroup.ExitNotify()

	for {
		select {
		case <-exitNotify:
			return
		case reqData := <-this.ReqMsgCh:
			logs.Debug(reqData)
		case respData := <-this.RespMsgCh:
			this.writeNetData(respData)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

/*
func init() {
	var rlim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim); err != nil {
		panic(err.Error())
	}
	rlim.Cur = 1000000
	rlim.Max = 1000000
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim); err != nil {
		panic(err.Error())
	}
}
*/

func main() {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9100")
	if err != nil {
		logs.Debug(err)
		return
	}
	client, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logs.Debug(err)
		return
	}
	logs.Debug("connect success")

	g_client.Conn = client

	go g_client.Serve()

	var str string
	obj := json.New()
	for {
		fmt.Scanf("%s", &str)
		obj.Set("content", str)
		g_client.ReqMsgCh <- obj
	}
}
