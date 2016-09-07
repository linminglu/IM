package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamQueryReq struct {
	TeamIdList []uint64 `json:"teamidlist"`
}

type TeamQueryResp struct {
	TeamInfolist []common.TeamInfo `json:"teamInfolist"`
}

func (h *TeamHandler) Query(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TeamQueryReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	var infolist []common.TeamInfo

	for _, val := range req.TeamIdList {
		teamInfo := &common.TeamInfo{TeamId: val}
		if RedisQueryInfo(teamInfo) == 0 {
			infolist = append(infolist, *teamInfo)
		}
	}

	if len(infolist) == 0 {
		return []byte(`{"teamInfolist":[]}`), 0
	}

	resp := TeamQueryResp{infolist}
	respBuf, err := json.Marshal(resp)

	return []byte(respBuf), 0
}
