package rest_server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/donnie4w/go-logger/logger"
	"github.com/vmihailenco/redis/v2"
	"sirendaou.com/duserver/common"
)

type OfflineMsgReq struct {
	Uid   uint64 `json:"uid,omitempty"`
	Count int    `json:"count,omitempty"`
}

type OfflineMsgResp struct {
	State int                   `json:"state"`
	Msg   string                `json:"msg"`
	Msgs  []*common.UserMsgItem `json:"msgs,omitempty"`
}

func (h *Handler) OffLineMsg(res http.ResponseWriter, req *http.Request) {
	restResp := OfflineMsgResp{State: 0, Msg: "ok", Msgs: []*common.UserMsgItem{}}

	defer func() {
		result, err := json.Marshal(restResp)

		if err != nil {
			res.Write([]byte(`{"state":500,"msg":"server err"}`))
			logger.Error(err)
		} else {
			res.Write(result)
		}
	}()

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		restResp.State = 400
		restResp.Msg = "request json error"
		logger.Error(err)
		return
	}

	logger.Debug("rest api offline_msg req:", string(body))

	var msgReq OfflineMsgReq
	err = json.Unmarshal(body, &msgReq)

	if err != nil {
		restResp.State = 400
		restResp.Msg = "request json error"
		logger.Error(err)
		return
	}

	state, err := h.redisMgr.RedisMsgCacheGet(msgReq.Uid)

	if err != nil {
		restResp.State = 500
		restResp.Msg = "sys error"
		logger.Error(err)
		return
	}

	if state == "read" {
		logger.Debug(msgReq.Uid, " have no db msg")
		restResp.State = 500
		restResp.Msg = "have no msg"
		logger.Error(err)
		return
	}

	msgItem := common.UserMsgItem{}
	msgs, err := msgItem.GetUserMsg(msgReq.Uid)

	if err != nil {
		logger.Debug(msgReq.Uid, " have no db msg")
		restResp.State = 500
		restResp.Msg = "have no msg"
		logger.Error(err)
		return
	}

	restResp.Msgs = append(restResp.Msgs, msgs...)

	// 客户端是否确认
	for _, msg := range msgs {
		msg.DelUserMsg()
	}

	msgs = CacheGetUserMsg(msgReq.Uid)
	restResp.Msgs = append(restResp.Msgs, msgs...)
	logger.Debug("offline msg count:", len(msgs))

	h.redisMgr.RedisMsgCacheSet(msgReq.Uid, "read")
}

//==================================获得群离线消息===========================================
var g_msgredis *common.RedisManager = nil

func TeamMsgRedisInit(redisAddr string) int {
	g_msgredis = common.NewRedisManager(redisAddr)
	if g_msgredis == nil {
		return -1
	}
	return 0
}

func CacheGetMsgBuf(msgId uint64) []byte {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAMMSGBUF, msgId)
	val, err := rClient.Client.Get(userKey).Result()

	if err != nil {
		logger.Error("error:", err)
		return nil
	}

	return []byte(val)
}

func CacheGetMsgIds(uid uint64) []uint64 {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAMMSGID, uid)
	vals, err := rClient.Client.ZRange(userKey, 0, -1).Result()

	if err != nil {
		if err == redis.Nil {
			logger.Info("ZRange ", userKey, 0, -1, " not data")
		} else {
			logger.Info("ZRange ", userKey, 0, -1, "fail", err.Error())
		}

		return nil
	}

	uidList := make([]uint64, len(vals))
	for i, s := range vals {
		uidList[i], err = strconv.ParseUint(s, 10, 64)
	}

	return uidList
}

func CacheGetUserMsg(uid uint64) []*common.UserMsgItem {
	msgIdList := CacheGetMsgIds(uid)

	msgList := []*common.UserMsgItem{}
	msg := &common.TextTeamMsg{}
	for _, msgId := range msgIdList {
		msgBuf := CacheGetMsgBuf(msgId)
		ret, _, jsonBody, _ := common.DecPkgInnerBody(msgBuf)
		if ret != 0 {
			logger.Info("DecPkgInnerBody failed:", ret)
			continue
		}

		if err := json.Unmarshal(jsonBody, msg); err != nil {
			logger.Info("json.Unmarshal failed:", err)
			continue
		}

		userMsg := new(common.UserMsgItem)
		userMsg.CmdType = common.DU_PUSH_CMD_IM_TEAM_MSG
		userMsg.MsgId = msg.MsgId
		userMsg.FromUid = msg.FromUid
		userMsg.ToUid = msg.ToTeamId
		userMsg.Type = uint16(msg.MsgType)
		userMsg.Content = msg.MsgContent
		userMsg.SendTime = uint32(msg.SendTime)
		userMsg.ApnsText = msg.ApnsText
		userMsg.FBv = msg.FBv
		//UserMsg.
		msgList = append(msgList, userMsg)
	}

	return msgList
}
