package team

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamMsgSendReq struct {
	MsgContent string `json:"msgcontent"`
	TeamId     uint64 `json:"toteamid"`
	MsgType    int    `json:"msgtype"`
	ApnsText   string `json:"apnstext"`
	FBv        int    `json:"fbv"`
}

func (h TeamHandler) SendMsg2Mq(head common.PkgHead, req common.TextTeamMsg, tail *common.InnerPkgTail, isSave bool) {
	jsonBody, _ := json.Marshal(req)

	msgBuf := new(bytes.Buffer)

	head.PkgLen = common.SIZEOF_PKGHEAD + uint16(len(jsonBody))
	head.Seq = 0
	head.Sid = 0 // errcode
	head.Cmd = common.DU_PUSH_CMD_IM_TEAM_MSG

	binary.Write(msgBuf, binary.BigEndian, head)
	binary.Write(msgBuf, binary.BigEndian, jsonBody)
	binary.Write(msgBuf, binary.BigEndian, tail)

	logger.Debug("teamMsg send msg to", tail.ToUid, "msgid:", tail.MsgId)

	h.chatMsgCh <- msgBuf.Bytes()

	if isSave {
		CacheSetMsgBuf(req.MsgId, msgBuf.Bytes())
	}

	//h.msgsavech <- &comm.TeamMsgSaveItem{1, tail.Msgid, tail.Touid, req.Toteamid, msgbuf.Bytes(), uint32(req.SendTime)}
}

func (h *TeamHandler) SendMsg(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamMsgSendReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	teamInfo := &common.TeamInfo{Uid: tail.FromUid, TeamId: req.TeamId}

	ret, uidList := RedisQueryMembers(teamInfo)
	if ret != 0 {
		logger.Error("RedisQueryMembers fail")
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	logger.Info("RedisQueryMembers  ok, uidList size ", len(uidList))

	var teamMsg common.TextTeamMsg
	teamMsg.SendTime = int(time.Now().Unix())
	teamMsg.FromUid = head.Uid
	teamMsg.MsgContent = req.MsgContent
	teamMsg.MsgType = req.MsgType
	teamMsg.ToTeamId = req.TeamId
	teamMsg.ApnsText = req.ApnsText

	if req.FBv > 0 {
		teamMsg.FBv = req.FBv
	} else {
		vk := fmt.Sprintf("BV_%d", head.Uid)
		bv, _ := h.userRedisMgr.RedisGet(vk)
		teamMsg.FBv, _ = strconv.Atoi(bv)
	}

	logger.Info("teamMsg.FBv :", teamMsg.FBv)

	MsgId := common.GetTeamMsgId()

	isSave := 0
	for i, uid := range uidList {
		if uint64(uid) == head.Uid {
			logger.Info("not send to self ", uid)
			uidList[i] = 0
			continue
		}
		tail.ToUid = uint64(uid)
		tail.FromUid = uint64(head.Uid)
		tail.MsgId = MsgId
		teamMsg.MsgId = tail.MsgId
		if isSave == 0 {
			h.SendMsg2Mq(head, teamMsg, tail, true)
		} else {
			isSave = 1
			h.SendMsg2Mq(head, teamMsg, tail, false)
		}
	}

	score := float64(teamMsg.SendTime)
	CacheAddMsgId(uidList, MsgId, score)

	results := fmt.Sprintf(`{"msgid":%d,"sendtime":%d}`, MsgId, teamMsg.SendTime)

	return []byte(results), 0
}
