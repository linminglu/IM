package db

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type GetCSListReq struct {
	Appkey  string   `json:"appkey,omitempty"`
	Uidlist []uint64 `json:"uidlist,omitempty"`
}

type GetCSListResp struct {
	Cslist []common.CSInfo `json:"cslist,omitempty"`
}

func (h *DBHandler) GetCSList(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req GetCSListReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if len(req.Appkey) != 24 {
		logger.Error("len(req.Appkey) != 24")
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	userinf := common.UserInfo{}
	n, csinfs := userinf.DBGetCSListByUidList(req.Uidlist)

	if n != 0 {
		logger.Error("DBGetCSList have no record")
		return []byte(""), common.ERR_CODE_SYS
	}

	resp := GetCSListResp{csinfs}

	respbuf, err := json.Marshal(resp)

	if err != nil {
		logger.Error("Marshal error:", err)
		return []byte(""), common.ERR_CODE_SYS
	}

	return respbuf, 0
}
