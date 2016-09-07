package kefu

import (
	"fmt"
	"time"

	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type SaveMsg struct {
	FromNick   string `json:"from_nick,omitempty"`
	FromCid    string `json:"from_cid,omitempty"`
	ToNick     string `json:"to_nick,omitempty"`
	ToCid      string `json:"to_cid,omitempty"`
	SendTime   int    `json:"send_time,omitempty"`
	MsgType    int    `json:"msg_type,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type RecvReq struct {
	To_user  int64  `json:"to_user,omitempty"`
	Msg_type string `json:"msg_type,omitempty"`
	Content  string `json:"content,omitempty"`
}

func (h *Handler) Recv(ctx *web.Context, val string) {
	logger.Debug("recv head:", ctx.Request.Header)
	logger.Debug("recv Form:", ctx.Request.Form)
	logger.Debug("recv PostForm:", ctx.Request.PostForm)

	retStr := ""

	defer func() {
		logger.Info("return:", retStr)
		ctx.Write([]byte(retStr))
	}()

	retCookie, ok := ctx.Request.Cookie("JSESSIONID")
	logger.Info("cookie:", retCookie, ok)

	logger.Info(h.Check(retCookie.Value))
	app_key, errcode, uid, account := h.Check(retCookie.Value)
	logger.Debug("error:", errcode)

	if errcode != 0 {
		switch errcode {
		case 1:
			retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s"}`, "")
		case 2:
			retStr = fmt.Sprintf(`{"code":1002,"err_msg":"%s"}`, "")
		default:
			retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "")

		}
		return
	}

	mqkey := fmt.Sprintf("csmsg_%d", uid)
	tnow := time.Now().Unix()
	Cid := ""
	Msg := ""
	STime := ""
	msgbase64 := ""
	Send_time := 0
	MsgType := 1
	n := 0
	for n < 30 {
		time.Sleep(time.Second * 1)
		val := h.MsgRedis.RedisRPop(mqkey)

		if len(val) > 10 {
			logger.Info("rpop mqkey", val)
			z := strings.SplitN(val, "|", 4)

			if len(z) == 4 {
				Cid = z[0]
				STime = z[1]
				Msg = z[2]
				MsgType, _ = strconv.Atoi(z[3])
				msgbase64 = base64.StdEncoding.EncodeToString([]byte(Msg))

				Send_time, _ = strconv.Atoi(STime)
				if Send_time+7200 < int(tnow) {
					logger.Info("msg time out")
				} else {
					break
				}
			}
		}
		n++
	}

	if len(Msg) > 1 {
		retStr = fmt.Sprintf(`{"code":0, "data":{ "total":1,"list":[{"type":3,"content":{"msg_id":%d,"msg_type":%d,"from_user":"%s","from_nick":"%s","status":1,"content":"%s", "send_time":%d}}]}}`, tnow, MsgType, Cid, Cid, msgbase64, Send_time)

		stMsg := SaveMsg{Cid, Cid, account, account, Send_time, MsgType, msgbase64}
		msgbuf, err := json.Marshal(stMsg)
		if err == nil {
			common.DBCSInsertMsg(Cid, account, string(msgbuf[:]), app_key, int(tnow))
		} else {
			logger.Info("Marshal ", stMsg, " fail")
		}
	} else {
		retStr = fmt.Sprintf(`{"code":1005,"err_msg":"%s"}`, "have no message")
	}

	return
}
