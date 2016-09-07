package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamDeleteReq struct {
	TeamId uint64 `json:"teamid,omitempty"`
}

func (h *TeamHandler) Delete(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamDeleteReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	logger.Info("uid :", head.Uid, "  remove teamid:", req.TeamId)

	if (req.TeamId & common.UID_FLAG) != (head.Uid & common.UID_FLAG) {
		logger.Error("teamid is error")
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	teamInfo := &common.TeamInfo{Uid: head.Uid, TeamId: req.TeamId}

	teamInfo.DBDelete()
	CacheDelete(teamInfo)

	return []byte("{}"), 0
}
