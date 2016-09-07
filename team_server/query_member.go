package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamQueryMemberReq struct {
	TeamId uint64 `json:"teamid"`
}

type TeamQueryMemberResp struct {
	Members []uint64 `json:"members"`
}

func (h *TeamHandler) QueryMember(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamQueryMemberReq

	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	logger.Info("uid %d query teamid %d 's info", head.Uid, req.TeamId)

	teamInfo := &common.TeamInfo{Uid: head.Uid, TeamId: req.TeamId}

	errCode, uidList := RedisQueryMembers(teamInfo)
	if errCode != 0 {
//		return nil, common.ERR_CODE_SYS
		uidList, err = teamInfo.DBGetMembers()
		if err != nil {
			return nil, common.ERR_CODE_SYS
		} else {
			for uid := range uidList {
				CacheAddMember(teamInfo, uint64(uid))
			}
		}
	}

	resp := TeamQueryMemberResp{uidList}
	respBuf, err := json.Marshal(resp)

	return []byte(respBuf), 0
}
