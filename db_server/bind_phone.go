package db

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type BindPhoneReq struct {
	PhoneNum string `json:"phonenum,omitempty"`
}

func (h *DBHandler) BindPhone(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req BindPhoneReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if head.Uid == 0 || len(req.PhoneNum) == 0 {
		logger.Error("BindPhone error req.Uid=", head.Uid)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	userInfo, err := common.GetUserInfoByUid(head.Uid)
	if err != nil {
		logger.Error(err)
		return []byte(""), common.ERR_CODE_NO_USER
	}

	userInfo.PhoneNum = req.PhoneNum
	logger.Debug(userInfo.PhoneNum)

	if err := userInfo.DBInsertPhoneNum(); err != nil {
		logger.Error(err)
		return []byte(""), common.ERR_CODE_SYS
	}

	return []byte(""), 0
}
