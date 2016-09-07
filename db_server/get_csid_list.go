package db

import (
	"encoding/json"
	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type GetCSIdListReq struct {
	Appkey string `json:"appkey,omitempty"`
	V      int    `json:"v,omitempty"`
}

type GetCSIdListResp struct {
	Csidverlist []common.CSVerInfo `json:"csidverlist,omitempty"`
}

func (h *DBHandler) GetCSIdList(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req GetCSIdListReq
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

	_, csids := userinf.DBGetCSVerList(req.V)

	resp := GetCSIdListResp{csids}
	respbuf, err := json.Marshal(resp)

	if err != nil {
		logger.Error("Marshal error:", err)
		return []byte(""), common.ERR_CODE_SYS
	}

	return respbuf, 0
}
