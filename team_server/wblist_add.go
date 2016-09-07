package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type WBAddReq struct {
	Uid  uint64 `json:"uid,omitempty"`
	Type int    `json:"type,omitempty"`
}

func (h *TeamHandler) WBAdd(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req WBAddReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	logger.Info("WBAddReq fromuid", head.Uid, " add ", req.Uid, "type:", req.Type)

	if req.Type != 1 && req.Type != 2 {
		logger.Error("Type is error:", req.Type)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	teamInfo := &common.TeamInfo{Uid: head.Uid}

	teamInfo.DBWBAdd(req.Uid, req.Type)
	CacheAddWBMember(teamInfo, req.Uid, req.Type)

	//if add blacklist, need cancel the friendship
	if req.Type == 2 {
		teamInfo.DBWBDelete(req.Uid, 1)
		CacheRemoveWBMember(teamInfo, req.Uid, 1)

		teamInfo2 := &common.TeamInfo{Uid: req.Uid}
		teamInfo2.DBWBDelete(head.Uid, 1)
		CacheRemoveWBMember(teamInfo2, head.Uid, 1)
	}

	return []byte(""), 0
}
