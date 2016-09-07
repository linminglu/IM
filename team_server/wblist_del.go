package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type WBDeleteReq struct {
	Uid  uint64 `json:"uid,omitempty"`
	Type int    `json:"type,omitempty"`
}

func (h *TeamHandler) WBDelete(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req WBDeleteReq
	err := json.Unmarshal(jsonBody, &req)

	logger.Info("WBDelete fromuid", head.Uid, " del ", req.Uid, "type:", req.Type)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	if req.Type != 1 && req.Type != 2 {
		logger.Error("Type is error:", req.Type)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	teamInfo := &common.TeamInfo{Uid: head.Uid}

	teamInfo.DBWBDelete(req.Uid, req.Type)

	CacheRemoveWBMember(teamInfo, req.Uid, req.Type)

	return []byte(""), 0
}
