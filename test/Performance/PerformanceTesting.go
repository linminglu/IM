package main


import (
	"bytes"
	"encoding/json"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
	"flag"
	"io"
)

var (
	g_ListenAddr = flag.String("listenaddr", "0.0.0.0:8080", "redis mq server addr")
	g_start      = flag.Int("start", 1, "start account")
	g_end        = flag.Int("end", 1, "end account")
	g_appkey     = flag.String("appkey", "" , "appkey")
)

type PkgHead struct {
	Pkglen uint16
	Cmd    uint16
	Ver    uint16
	Seq    uint16
	Sid    uint32
	Uid    uint64
	Flag   uint32
}

type LoginResp struct {
	Uid        uint64 `json:"uid,omitempty"`
	Session_id uint32   `json:"sid,omitempty"`
}

func client(n int,  succch, quitch chan int) {
    
	fmt.Println("cleint ", n , "start")
	server := *g_ListenAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)

	isOK := false

	defer func() {
	    fmt.Println("cleint ", n , " end ,success ?", isOK)
		if isOK {
			succch <- n
		} else {
			quitch <- n
		}
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println("connect fail")
		return
	}

	defer conn.Close()

	reqBuf := new(bytes.Buffer)

	loginjson := fmt.Sprintf(`{"platform":"i","appkey":"%s", "cid":"nhtest_%07d","password":"453e41d218e071ccfb2d1c99ce23906a"}`, *g_appkey, n)
	head := PkgHead{}
	head.Pkglen = uint16(24 + len(loginjson))
	head.Cmd = 10101 
	head.Ver = 1
	head.Seq = 1
	head.Sid = 0
	head.Flag = 0
	binary.Write(reqBuf, binary.BigEndian, head)
	binary.Write(reqBuf, binary.BigEndian, []byte(loginjson))

	conn.Write(reqBuf.Bytes())

	//read resp
	nethead := make([]byte, 2)
	conn.SetDeadline(time.Now().Add(time.Second * 60))

	if nrR, e := io.ReadFull(conn, nethead[:2]); e != nil || nrR != 2 {
		fmt.Println("read resp timeout")
		return
	}
	headlen := binary.BigEndian.Uint16([]byte(nethead))

	fmt.Println( "heand len :", headlen)

	if headlen > 128 {
		fmt.Println("heand len :", headlen, " to large")
		return
	} else if headlen < 24 {
		fmt.Println("heand len :", headlen, " to small")
		return
	}

	respbuf := make([]byte, int(headlen))
	respbuf[0] = nethead[0]
	respbuf[1] = nethead[1]
	conn.SetDeadline(time.Now().Add(time.Second * 5))
	if nr, e := io.ReadFull(conn, respbuf[2:]); e != nil || nr+2 != len(respbuf) {

		fmt.Println("recv error:", e)
		return
	}

	p := bytes.NewReader(respbuf)

	err = binary.Read(p, binary.BigEndian, &head)
	if err != nil {
		fmt.Println("read pkghead fail:", err.Error())
		return
	}

	if head.Sid != 0 {
		fmt.Println("login errcode:", head.Sid)
		return 
	}

	retjson := string(respbuf[24:])
	fmt.Println(retjson)

	var resp LoginResp
	err = json.Unmarshal(respbuf[24:], &resp)

	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	isOK = true

	for i := 0; i < 50; i++ {

	    fmt.Println("start send hello to conn " , n)
		time.Sleep(time.Second * 600)
		head.Pkglen = 24
		head.Cmd = 10100
		head.Sid = resp.Session_id
		head.Uid = resp.Uid

	    reqBuf := new(bytes.Buffer)
		binary.Write(reqBuf, binary.BigEndian, head)
		err = conn.SetWriteDeadline(time.Now().Add(time.Second*10))
		l, err := conn.Write(reqBuf.Bytes())
		if err != nil {
		    fmt.Println(l, err)
	    }

	    conn.Read(respbuf) 
	}
}

func main() {
    flag.Parse()

	succch :=make(chan int , 1)
	failch :=make(chan int , 1)
	for i := *g_start ; i < *g_end; i++{
       go client(i , succch, failch)	
	   time.Sleep(time.Millisecond*30)
	}

	for{
        select {
	        case  n := <- succch :
			    fmt.Println(n, " success")	
	        case  n := <- failch:
			    fmt.Println(n, " fail")	
		}	
	}
	return
}
