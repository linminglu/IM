package team

import (
	"encoding/json"
	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamUpdateReq struct {
	TeamId   uint64 `json:"teamid"`
	CoreInfo string `json:"coreinfo"`
	ExInfo   string `json:"exinfo"`
	TeamName string `json:"teamname"`
}

func (h *TeamHandler) Update(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamUpdateReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	logger.Info("uid ", head.Uid, " update teamid", req.TeamId, "coreinfo ", req.CoreInfo, " ,exinfo:", req.ExInfo)

	if (head.Uid & common.UID_FLAG) != (req.TeamId & common.UID_FLAG) {
		logger.Error("uid ", head.Uid, " is not  tid", req.TeamId, " creator")
		return []byte(""), common.ERR_CODE_TEAM_PRI
	}

	teamInfo := &common.TeamInfo{Uid: head.Uid, TeamId: req.TeamId, TeamName: req.TeamName, CoreInfo: req.CoreInfo, ExInfo: req.ExInfo}
	teamInfo.DBSetInfo()
	CacheSetInfo(teamInfo)

	return []byte(""), 0
}
