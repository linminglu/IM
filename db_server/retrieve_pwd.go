package db

import (
	"encoding/json"
	//	"strconv"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type RetrievePwdReq struct {
	PhoneNum string `json:"phonenum,omitempty"`
	Pwd      string `json:"password,omitempty"`
}

//type RetrievePwdResp struct {
//	PhoneNum string `json:"phonenum,omitempty"`
//	Password string `json:"password,omitempty"`
//}

func (h *DBHandler) RetrievePwd(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	logger.Debug("RetrievePwd head=", head.ToString())

	var req RetrievePwdReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if req.PhoneNum == "" || req.Pwd == "" {
		logger.Error("nil arguments. Phonenum: ", req.PhoneNum, " Password: ", req.Pwd)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	//	phonenum_uid_key := "du_phonenum_uid@" + req.PhoneNum
	//	uidStr, code := h.UserRedis.RedisGet(phonenum_uid_key)
	//	if code != 0 || uidStr == "" {
	uid, err := common.GetUidByPhoneNum(req.PhoneNum)
	if err != nil {
		logger.Debug("user not exists Phonenum=", req.PhoneNum)
		return []byte(""), common.ERR_CODE_NO_USER
	}
	//	}

	retCode := common.DBUpdateUserInfo(uid, req.Pwd)
	if retCode != 0 {
		logger.Error("DBUpdateUserInfo err=", err.Error())
		return []byte(""), common.ERR_CODE_SYS
	}

	return []byte(""), 0
}
