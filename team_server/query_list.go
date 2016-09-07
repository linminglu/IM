package team

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type TeamQueryListReq struct {
	Uid int64 `json:"uid"`
}

type TeamQueryListResp struct {
	TeamIdList []int64 `json:"teamidlist"`
}

func (h *TeamHandler) QueryList(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	logger.Info("QueryList, req:", string(jsonBody[:]))

	var req TeamQueryListReq

	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERR_CODE_ERR_PKG
	}

	tail.FromUid = head.Uid

	teamInfo := &common.TeamInfo{Uid: head.Uid}

	errCode, tidList := RedisQueryList(teamInfo)
	if errCode != 0 {
		logger.Error("RedisQueryList fail , ret %d", errCode)
		return nil, common.ERR_CODE_SYS
	}

	if tidList == nil || len(tidList) < 1 {
		return []byte(`{"teamidlist":[]}`), 0
	}

	resp := TeamQueryListResp{tidList}
	respBuf, err := json.Marshal(resp)

	logger.Info("QueryList, resp:", string(respBuf[:]))

	return []byte(respBuf), 0
}

func (h *TeamHandler) QuerySysList(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	errCode, tidList := common.DBGetSysTeamIdList()
	if errCode != 0 {
		logger.Error("DBGetSysTeamIdList fail, ret %d", errCode)
		return nil, common.ERR_CODE_SYS
	}

	if tidList == nil || len(tidList) < 1 {
		return []byte(`{"teamidlist":[]}`), 0
	}

	resp := TeamQueryListResp{tidList}

	respBuf, err := json.Marshal(resp)
	if err != nil {
		return nil, common.ERR_CODE_SYS
	}

	logger.Info("QueryList, resp:", string(respBuf[:]))

	return []byte(respBuf), 0
}