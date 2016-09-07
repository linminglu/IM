package db

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"
	"sirendaou.com/duserver/common"
)

type AppCfgReq struct {
	Appkey string `json:"appkey,omitempty"`
}

type AppCfgResp struct {
	Kefu  string `json:"kefu,omitempty"`
	Csver int    `json:"csver,omitempty"`
}

func (h *DBHandler) AppCfg(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req AppCfgReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	userinf := common.UserInfo{}
	v := userinf.DBGetCSLastVer()

	resp := AppCfgResp{"", v}
	respbuf, err := json.Marshal(resp)

	return respbuf, 0
}
