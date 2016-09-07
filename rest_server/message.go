package rest_server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type Msg struct {
	Type   string `json:"type,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Action string `json:"action,omitempty"`
}

//type MsgExt struct {
//	Attr1 string `json:"attr1,omitempty"`
//	Attr2 string `json:"attr2,omitempty"`
//}

type SendReq struct {
	FromUser uint64   `json:"from"`
	ToUsers  []uint64 `json:"target"`
	Message  Msg      `json:"msg"`
	Ext      string   `json:"ext"`
}

func (req SendReq) genContent() string {
	req.ToUsers = nil
	result, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return ""
	}

	//	logger.Debug("genContent================", string(result))

	return string(result)
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
	msg.CmdType = common.DU_PUSH_CMD_IM_SYSTEM_MSG
	// 透传消息
	if req.Message.Type == "cmd" {
		msg.Content = req.Message.Action
	} else {
		msg.Content = req.genContent()
	}
	msg.ApnsText = ""
	msg.FType = 1 //来自客服

	jsonBody, _ := json.Marshal(msg)

	msgBuf := new(bytes.Buffer)

	head.PkgLen = common.SIZEOF_PKGHEAD + uint16(len(jsonBody))
	head.Seq = 0
	head.Sid = 0 // errcode
	head.Cmd = common.DU_PUSH_CMD_IM_SYSTEM_MSG

	//	logger.Debug("---------------jsonBody===", string(jsonBody))

	binary.Write(msgBuf, binary.BigEndian, head)
	binary.Write(msgBuf, binary.BigEndian, jsonBody)
	binary.Write(msgBuf, binary.BigEndian, tail)

	err := h.producer.Publish(*g_Db2MsgTopic, msgBuf.Bytes())
	if err != nil {
		logger.Error("nsq write to ", *g_Db2MsgTopic, "err:", err.Error())
	} else {
		logger.Debug("nsq write to ", *g_Db2MsgTopic, " success")
	}

	if err := msg.SaveUserMsg(string(msgBuf.Bytes())); err != nil {
		logger.Error(err)
	}

	return 0
}

//req body: {"from":1001,"target":[11175,100042,100043],"msg":{"type":"txt", "action":"action", "msg":"this is sys msg"},"ext":{"attr1":"attr1","attr2":"attr2"}}
func (h *Handler) SendMessage(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		res.Write([]byte(`{"state":400,"msg":"request json error"}`))
		return
	}

	logger.Debug("SendMessage request body==== \n", string(body))

	var sendReq SendReq
	err = json.Unmarshal(body, &sendReq)
	if err != nil {
		logger.Error(err)
		res.Write([]byte(`{"state":400,"msg":"request json error"}`))
		return
	}

	for _, uid := range sendReq.ToUsers {
		h.SendMsg2Mq(sendReq.FromUser, uid, sendReq)
	}

	res.Write([]byte(`{"state":0,"msg":"ok"}`))

	return
}
