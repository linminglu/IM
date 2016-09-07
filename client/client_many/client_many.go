package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"syscall"
	"time"

	"sirendaou.com/duserver/common"
	db_server "sirendaou.com/duserver/db_server"
)

var (
	g_client *net.TCPConn = nil
	g_uid    int64        = 0
	g_sid    int          = 0
)

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
func sendPkg(client *net.TCPConn, head *common.PkgHead, body interface{}) error {
	buf := new(bytes.Buffer)
	bodyJson, err := json.Marshal(body)
	if err != nil {
		fmt.Println("json.Marshal error:")
		return err
	}
	head.PkgLen = common.SIZEOF_PKGHEAD + uint16(len(bodyJson))
	binary.Write(buf, binary.BigEndian, head)
	binary.Write(buf, binary.BigEndian, bodyJson)
	if _, err := client.Write(buf.Bytes()); err != nil {
		fmt.Println("client.Write error")
		return err
	}
	//fmt.Println("send head:", *head)
	//fmt.Println("send body:", body)
	return nil
}

func recvPkg(client *net.TCPConn, head *common.PkgHead, body interface{}) error {
	client.SetDeadline(time.Now().Add(time.Hour))
	pkgLenSlice := make([]byte, 2)

	if _, err := io.ReadFull(client, pkgLenSlice); err != nil {
		fmt.Println("io.ReadFull error", err)
		return err
	}
	pkgLenInt := binary.BigEndian.Uint16(pkgLenSlice)
	respBuf := make([]byte, pkgLenInt)
	respBuf[0] = pkgLenSlice[0]
	respBuf[1] = pkgLenSlice[1]

	if _, err := io.ReadFull(client, respBuf[2:]); err != nil {
		fmt.Println("io.ReadFull error", err)
		return err
	}

	respReader := bytes.NewReader(respBuf)
	if err := binary.Read(respReader, binary.BigEndian, head); err != nil {
		fmt.Println("binary.Read error:")
		return err
	}
	jsonStr := make([]byte, head.PkgLen-common.SIZEOF_PKGHEAD)
	if err := binary.Read(respReader, binary.BigEndian, jsonStr); err != nil {
		fmt.Println("binary.Read error:")
		return err
	}
	if len(jsonStr) != 0 && body != nil {
		if err := json.Unmarshal(jsonStr, body); err != nil {
			fmt.Println("json.Unmarshal error:")
			return err
		}
	}
	/*
		fmt.Println("recv head:", head)
		if body == nil {
			fmt.Println("recv body:", string(jsonStr))
		} else {
			fmt.Println("recv body:", body)
		}
	*/
	return nil
}

func hello(client *net.TCPConn, sid int, uid int64) error {
	//fmt.Println("===============================hello==================================")
	buf := new(bytes.Buffer)
	pkgHead := &common.PkgHead{
		PkgLen: common.SIZEOF_PKGHEAD,
		Cmd:    common.DU_CMD_USER_HELLO,
		Ver:    1,
		Seq:    0,
		Sid:    uint32(sid),
		Uid:    uint64(uid),
		Flag:   0,
	}
	binary.Write(buf, binary.BigEndian, pkgHead)
	client.SetDeadline(time.Now().Add(time.Hour))
	if _, err := client.Write(buf.Bytes()); err != nil {
		return err
	}
	//fmt.Println("send:", pkgHead)
	pkgLenSlice := make([]byte, 2)

	if _, err := io.ReadFull(client, pkgLenSlice); err != nil {
		return err
	}
	pkgLenInt := binary.BigEndian.Uint16(pkgLenSlice)
	respBuf := make([]byte, pkgLenInt)
	respBuf[0] = pkgLenSlice[0]
	respBuf[1] = pkgLenSlice[1]

	if _, err := io.ReadFull(client, respBuf[2:]); err != nil {
		return err
	}
	respReader := bytes.NewReader(respBuf)
	respHead := &common.PkgHead{}
	if err := binary.Read(respReader, binary.BigEndian, respHead); err != nil {
		return err
	}
	//fmt.Println("recv head:", respHead)
	return nil
}

func register(client *net.TCPConn) {
	fmt.Println("===============================register==================================")
	pkgHead := &common.PkgHead{
		PkgLen: common.SIZEOF_PKGHEAD,
		Cmd:    common.DU_CMD_USER_REGISTER,
		Ver:    1,
	}
	regReq := &db_server.RegReq{
		Platform: "a",
		DeviceId: "test123",
		PhoneNum: "18702759796",
		Password: "test",
		Login:    0,
	}
	if err := sendPkg(client, pkgHead, regReq); err != nil {
		fmt.Println("sendPkg error")
		return
	}
	respHead := &common.PkgHead{}
	respStruct := &db_server.RegResp{}
	if err := recvPkg(client, respHead, respStruct); err != nil {
		fmt.Println("recvPkg error")
		return
	}
	g_uid = respStruct.Uid
	g_sid = respStruct.Sid
	return
}

func login(client *net.TCPConn, platform string, uid uint64, passwd string) (int, int64, error) {
	//fmt.Println("===============================login==================================")
	pkgHead := &common.PkgHead{}
	pkgHead.Cmd = common.DU_CMD_USER_LOGIN
	regReq := &db_server.LoginReq{
		Platform: platform,
		Uid:      uid,
		Password: passwd,
	}
	if err := sendPkg(client, pkgHead, regReq); err != nil {
		return 0, 0, err
	}
	respHead := &common.PkgHead{}
	respStruct := &db_server.LoginResp{}
	if err := recvPkg(client, respHead, respStruct); err != nil {
		return 0, 0, err
	}
	g_uid = respStruct.Uid
	g_sid = respStruct.Sid
	return respStruct.Sid, respStruct.Uid, nil
}

type UserMsgReq struct {
	MsgContent string `json:"msgcontent,omitempty"`
	ToUid      uint64 `json:"touid,omitempty"`
	MsgType    uint16 `json:"msgtype,omitempty"`
	ApnsText   string `json:"apnstext,omitempty"`
	FBv        int    `json:"frombv,omitempty"`
	ExtraData  string `json:"extraData,omitempty"`
}

func sendMsg() {
	fmt.Println("===============================sendMsg==================================")
	pkgHead := &common.PkgHead{}
	pkgHead.Cmd = common.DU_CMD_IM_SEND_USER_MSG
	pkgHead.Uid = uint64(g_uid)
	pkgHead.Sid = uint32(g_sid)
	regReq := &UserMsgReq{
		MsgContent: "test" + fmt.Sprint(rand.Int31()),
		ToUid:      100057,
		ApnsText:   "test",
	}
	if err := sendPkg(g_client, pkgHead, regReq); err != nil {
		fmt.Println("sendPkg error")
		return
	}
	respHead := &common.PkgHead{}
	if err := recvPkg(g_client, respHead, nil); err != nil {
		fmt.Println("recvPkg error")
		return
	}
	return
}

type OfflineMsgReq struct {
	Uid   uint64 `json:"uid,omitempty"`
	Count int    `json:"count,omitempty"`
}

type OfflineMsgResp struct {
	State int                   `json:"state"`
	Msg   string                `json:"msg"`
	Msgs  []*common.UserMsgItem `json:"msgs,omitempty"`
}

func offlineMsg() {
	fmt.Println("===============================offlineMsg==================================")
	url := "http://127.0.0.1:3001/rest/offline_msg/list"
	//url := "http://192.168.20.51:3001/rest/offline_msg/list"
	req := &OfflineMsgReq{
		Uid:   100057,
		Count: 10,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}
	fmt.Println("content type:", http.DetectContentType(reqBody))
	resp, err := http.Post(url, http.DetectContentType(reqBody), bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println("http.Post error:", err)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("ioutil.ReadAll error:", err)
		return
	}
	respStruct := &OfflineMsgResp{}
	if err := json.Unmarshal(respBody, respStruct); err != nil {
		fmt.Println("json.Unmarshal error:", err)
		return
	}
	fmt.Println("recv msg", respStruct)
	for i, content := range respStruct.Msgs {
		fmt.Println("recv msg content ", i, ":", content)
	}
}

func main() {
	fmt.Println("===============================main==================================")
	for i := 0; i < 1000; i++ {
		go func(id int) {
		retry:
			//tcpAddr, err := net.ResolveTCPAddr("tcp", "192.168.20.51:9100")
			tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9100")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			client, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Println("connect_failed:", id, err)
				return
			}

			if err := client.SetKeepAlive(true); err != nil {
				fmt.Println("SetKeepAlive failed:", err)
				client.Close()
				goto retry
			}

			fmt.Println("connect_success:", id, client.LocalAddr())
			sid, uid, err := login(client, "a", 100056, "test")
			if err != nil {
				fmt.Println("login_failed", id, err)
				client.Close()
				return
			}
			fmt.Println("login_success:", id)
			for {
				if err := hello(client, sid, uid); err != nil {
					fmt.Println("hello_failed:", err, id)
					break
				}
				time.Sleep(time.Second)
			}

			client.Close()
		}(i)
	}
	fmt.Println("client exit")
	time.Sleep(time.Second * 120)
}
