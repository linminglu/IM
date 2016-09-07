package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamCreateReq struct {
	TeamName string `json:"teamname,omitempty"`
	TeamType int    `json:"teamtype,omitempty"` // 0-用户建的群组  1-系统预设群组
}

type TeamCreateResp struct {
	TeamId   uint64 `json:"teamid,omitempty"`
	MaxCount int    `json:"maxcount,omitempty"`
}

func (h *TeamHandler) Create(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamCreateReq

	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	logger.Info("uid ", head.Uid, " create team name", req.TeamName, " type ", req.TeamType)

	teamInfo := &common.TeamInfo{Uid: head.Uid, TeamType: req.TeamType, TeamName: req.TeamName, MaxCount: common.MAX_MEMBER_NUM_TRAM}

	tid := teamInfo.DBGetNewTeamID()
	logger.Debug("tid=",tid)
	if tid == 0 {
		return nil, common.ERR_CODE_SYS
	} else if tid == 1 {
		return nil, common.ERROR_TOUCH_TOP
	}

	teamInfo.TeamId = tid

	teamInfo.DBCreate()
	CacheCreate(teamInfo)

	teamInfo.DBAddMember(head.Uid)
	CacheAddMember(teamInfo, head.Uid)

	resp := TeamCreateResp{tid, common.MAX_MEMBER_NUM_TRAM}
	respBuf, err := json.Marshal(resp)
	return respBuf, 0
}
