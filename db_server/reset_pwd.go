package db

import (
	"encoding/json"

	"strconv"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type ResetPwdReq struct {
	OldPwd string `json:"oldpsw,omitempty"`
	NewPwd string `json:"newpsw,omitempty"`
}

func (h *DBHandler) ResetPwd(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	logger.Debug("ResetPwd()")
	var req ResetPwdReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if len(req.OldPwd) != 32 || len(req.NewPwd) != 32 {
		logger.Error("ResetPwd error")
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	tail.FromUid = uint64(head.Uid)

	v := strconv.FormatUint(uint64(head.Uid), 10)
	passKey := "Pwd_" + v
	passwd, _ := h.UserRedis.RedisGet(passKey)

	logger.Info("uid ", head.Uid, " input passwd", req.OldPwd, " right passwd:", passwd)

	//	if retCode != 0 || passwd == "" {
	//		TODO select from db
	//	}

	if passwd != req.OldPwd {
		logger.Info("passwd != req.OldPwd")
		return []byte(""), common.ERR_CODE_PASSWD
	}

	h.UserRedis.RedisSet(passKey, req.NewPwd)

	retCode := common.DBUpdateUserInfo(head.Uid, req.NewPwd)
	if retCode != 0 {
		return []byte(""), common.ERR_CODE_SYS
	}

	return []byte(""), 0
}
