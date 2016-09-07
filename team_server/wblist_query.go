package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type WBQueryReq struct {
	Type int `json:"type,omitempty"`
}

type WBQueryResp struct {
	UidList []uint64 `json:"uidList"`
}

func (h *TeamHandler) WBQuery(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req WBDeleteReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	logger.Info("WBQuery fromuid", head.Uid, "type:", req.Type)

	if req.Type != 1 && req.Type != 2 {
		logger.Error("Type is error:", req.Type)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	teamInfo := &common.TeamInfo{Uid: head.Uid}
	errCode, uidList := RedisQueryWBMembers(teamInfo, req.Type)
	if errCode != 0 {
		uidList, err = teamInfo.DBGetFriendUids(head.Uid)
		if err != nil {
			return nil, common.ERR_CODE_SYS
		} else {
			for uid := range uidList {
				CacheAddWBMember(teamInfo, uint64(uid), 1)
			}
		}
	}

	resp := WBQueryResp{uidList}
	respBuf, err := json.Marshal(resp)

	return []byte(respBuf), 0
}
