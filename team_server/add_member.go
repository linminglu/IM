package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamAddReq struct {
	TeamId uint64 `json:"teamid"`
	Uids   []uint64 `json:"uids"`
}

func (h *TeamHandler) AddMembers(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamAddReq

	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	} else if req.TeamId == 0 || len(req.Uids) == 0 {
		logger.Error("req.TeamId == 0 || req.Uid == 0")
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	if !RedisIsMembers(req.TeamId, head.Uid) {
		logger.Error("uid ", head.Uid, " is not  tid", req.TeamId, " creator")
		return []byte(""), common.ERR_CODE_TEAM_PRI
	}

	teamInfo := &common.TeamInfo{TeamId: req.TeamId}

	num := CacheScardMember(teamInfo)
	if num >= common.MAX_MEMBER_NUM_TRAM {
		return []byte(""), common.ERR_CODE_TRAM_MAXNUM
	}

	for _, uid := range req.Uids {
		teamInfo.DBAddMember(uid)
		CacheAddMember(teamInfo, uid)
	}

	return []byte(""), 0
}

