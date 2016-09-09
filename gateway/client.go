package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/astaxie/beego/logs"

	"IM/common/errors"
	"IM/common/json"
	sync_ "IM/common/sync"
)

type Client struct {
	Conn net.Conn
	Sid  int

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

func (this *Client) writeNetData(data *json.Json) *errors.ErrorT {
	buf := new(bytes.Buffer)

	body, err := json.Marshal(data)
	if err != nil {
		return errors.New(err.Error(), data.ToString())
	}

	dataSize := len(body)
	binary.Write(buf, binary.BigEndian, dataSize)
	binary.Write(buf, binary.BigEndian, body)

	this.Conn.SetDeadline(time.Now().Add(NET_WRITE_DEADLINE * time.Second))
	_, err = this.Conn.Write(buf.Bytes())

	if err != nil {
		return errors.New(err.Error())
	}
	return nil
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

	// 删除客户端
	g_gateway.clientMutex.Lock()
	delete(g_gateway.clientMap, this.Sid)
	g_gateway.clientMutex.Unlock()
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
			this.RespMsgCh <- reqData
		case respData := <-this.RespMsgCh:
			if err := this.writeNetData(respData); err != nil {
				logs.Debug(err.Error())
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
