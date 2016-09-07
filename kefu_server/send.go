package kefu

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type SendReq struct {
	ToUser  string `json:"to_user,omitempty"`
	MsgType uint16 `json:"msg_type,omitempty"`
	Content string `json:"content,omitempty"`
//	ExtraData string `json:"extraData,omitempty"`
}

func (h *Handler) SendMsg2Mq(from, to uint64, req SendReq) int {
	var msg common.UserMsgItem

	head := common.PkgHead{}
	tail := common.InnerPkgTail{}

	msg.SendTime = uint32(time.Now().Unix())
	msg.MsgId = common.GetChatMsgId()
	msg.FromUid = from
	msg.ToUid = to
	tail.ToUid = to
	tail.MsgId = msg.MsgId
	msg.Content = req.Content
	msg.Type = req.MsgType
	msg.ApnsText = ""
	msg.FType = 1 //来自客服
//	msg.ExtraData = req.ExtraData

	jsonBody, _ := json.Marshal(msg)

	msgBuf := new(bytes.Buffer)

	head.PkgLen = common.SIZEOF_PKGHEAD + uint16(len(jsonBody))
	head.Seq = 0
	head.Sid = 0 // errcode
	head.Cmd = common.DU_PUSH_CMD_IM_USER_MSG

	binary.Write(msgBuf, binary.BigEndian, head)
	binary.Write(msgBuf, binary.BigEndian, jsonBody)
	binary.Write(msgBuf, binary.BigEndian, tail)

	err := h.Producer.Publish(*g_Db2MsgTopic, msgBuf.Bytes())

	if err != nil {
		logger.Error("nsq write resp to ", *g_Db2MsgTopic, "err:", err.Error())
	} else {
		logger.Debug("nsq write resp to ", *g_Db2MsgTopic, " success")
	}

	if err := msg.SaveUserMsg(string(msgBuf.Bytes())); err != nil {
		logger.Error(err)
	}

	return 0
}

func (h *Handler) Send(ctx *web.Context, val string) {
	logger.Debug("Send head:", ctx.Request.Header)
	logger.Debug("Send Form:", ctx.Request.Form)
	logger.Debug("Send PostForm:", ctx.Request.PostForm)

	retStr := ""

	defer func() {
		ctx.Write([]byte(retStr))
	}()

	reqBuf := make([]byte, 1024)
	strLen, err := ctx.Request.Body.Read(reqBuf)

	logger.Debug("req:", string(reqBuf[:]))

	var req SendReq
	err = json.Unmarshal(reqBuf[0:strLen], &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "json body error")
		return
	}

	retCookie, ok := ctx.Request.Cookie("JSESSIONID")
	logger.Info("cookie:", retCookie.Value, ok)

	app_key, errcode, fromUid, account := h.Check(retCookie.Value)
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

	//	sKey := app_key + "_" + req.ToUser
	key := req.ToUser
	sUid, errcode := h.UserRedis.RedisGet(key)
	if errcode != 0 {
		retStr = fmt.Sprintf(`{"code":4001,"err_msg":""}`)
		return
	}

	toUid, _ := strconv.ParseUint(sUid, 10, 64)

	h.SendMsg2Mq(fromUid, toUid, req)
	retStr = fmt.Sprintf(`{"code":0,"msg_id":"%d"}`, time.Now().Unix())

	msgBase64 := base64.StdEncoding.EncodeToString([]byte(req.Content))
	sendTime := time.Now().Unix()
	stMsg := SaveMsg{account, account, req.ToUser, req.ToUser, int(sendTime), int(req.MsgType), msgBase64}
	msgBuf, err := json.Marshal(stMsg)
	if err == nil {
		common.DBCSInsertMsg(account, req.ToUser, string(msgBuf[:]), app_key, int(sendTime))
	} else {
		logger.Info("Marshal ", stMsg, " fail")
	}

	return
}
