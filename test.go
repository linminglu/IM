package main

import (
//	"time"
	"fmt"
//	"reflect"
//	"sirendaou.com/duserver/common"

	"github.com/bitly/go-simplejson"
)

type Html []interface{}

func main() {
	reqBuf := []byte(`[{"username":100, "password":"p1"}, {"username":100, "password":"p2"}]`)
//	reqBuf := []byte{"[{'username':'u1', 'password':'p1'}]"}

	retStr := ""

	js, err := simplejson.NewJson(reqBuf)
	if err != nil {
		retStr = "400"
		return
	}

	arr, err := js.Array()
	if err != nil || len(arr) > 60 {
		retStr = "400"
		return
	}

	strSql := "insert into t_user_info (id, uid, reg_date, update_date, password) values "
	for i := 0; i < len(arr); i++ {
		username, _ := js.GetIndex(i).Get("username").Uint64()
		password, _ := js.GetIndex(i).Get("password").String()
		s := fmt.Sprintf("(0 , %d, now(), now(), '%s'), ", username, password)
		strSql += s
	}

	strSql = strSql[:len(strSql)-2]

	fmt.Println("strSql=", strSql)
	fmt.Println("retStr=", retStr)
	fmt.Println("len(arr)=", len(arr))

//	logger.Debug("strSql=", strSql)

//	err = common.ProcExec(strSql)
//	if err != nil || len(arr) > 60 {
//		retStr = "400"
//		return
//	}

//	var val interface{} = "good"
////	val := "good"
//	fmt.Println(val.(string))
//
//	html := make(Html, 6)
//	html[0] = "div"
//	html[1] = "span"
//	html[2] = []byte("script")
//	html[3] = "style"
//	html[4] = "head"
//	html[5] = 100
//	for index, element := range html {
//		switch value := element.(type) {
//			case string:
//			fmt.Printf("html[%d] is a string and its value is %s\n", index, value)
//			case []byte:
//			fmt.Printf("html[%d] is a []byte and its value is %s\n", index, string(value))
//			case int:
//			fmt.Printf("html[%d] is a int and its value is %d\n", index, value)
//			default:
//			fmt.Printf("unknown type\n")
//		}
//	}

	//	t := time.Now()
	//	fmt.Println(t)
	//	t1 := time.Now().Unix()
	//	fmt.Println(t1)
	//	fmt.Println(reflect.TypeOf(t1))
	//	result := t.Format("2006-01-02 15:04:05")
	//	fmt.Println(result)

	//	redis_addr := "127.0.0.1:6379"
	//
	//	redisMgr := common.NewRedisManager(redis_addr)
	//	if redisMgr == nil {
	//		fmt.Println("connect redis ",redis_addr, "fail")
	//		return
	//	}

	//	type InnerPkgTail struct {
	//		ConnIP   int64
	//		ConnPort uint32
	//		FromUid  uint64
	//		ToUid    uint64
	//		Sid      uint32
	//		MsgId    uint64
	//	}

	//	tail := &common.InnerPkgTail {
	//		ConnIP:101010,
	//		ConnPort:9100,
	//		FromUid:10000,
	//		ToUid:20000,
	//		Sid:10000,
	//		MsgId:10101010,
	//	}
	//
	//	if err := redisMgr.RedisStatCacheSet(tail); err != nil {
	//		fmt.Println("err=",err)
	//	}
	//
	//	tailCache, err := redisMgr.RedisStatCacheGet(10000)
	//	if err != nil {
	//		fmt.Println("err=",err)
	//	}
	//
	//	fmt.Println("tailCache",tailCache)

}
