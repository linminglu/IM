package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"
	"sirendaou.com/duserver/common"
)

type MsgRecved struct {
	MsgId  uint64 `json:"msgid,omitempty"`
	TeamId uint64 `json:"teamid,omitempty"`
}

func (h TeamHandler) RecvedMsg(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req MsgRecved
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	logger.Info("revmove teamMsg touid:", head.Uid, " msgid:", req.MsgId)

	CacheRemMsgId(head.Uid, req.MsgId)

	return []byte(""), 0
}
