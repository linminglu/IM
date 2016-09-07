package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamRemReq struct {
	TeamId uint64 `json:"teamid"`
	Uid    uint64 `json:"uid"`
}

func (h *TeamHandler) RemoveMember(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamRemReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	} else if req.TeamId == 0 || req.Uid == 0 {
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	logger.Info("uid ", head.Uid, " remove from team id", req.TeamId, " a member %d) ", req.Uid)

	teamInfo := &common.TeamInfo{TeamId: req.TeamId}

	teamInfo.DBRemoveMember(req.Uid)
	CacheRemoveMember(teamInfo, req.Uid)

	return []byte(""), 0
}
