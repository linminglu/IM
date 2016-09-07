package db

import (
	"encoding/json"
	"fmt"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type SetSetupIdReq struct {
	SetupId uint64 `json:"setupid,omitempty"`
}

func (h *DBHandler) SetSetupId(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {

	logger.Info(head.Uid, "set setupid ", string(jsonBody[0:]))

	var req SetSetupIdReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	logger.Info(head.Uid, "SetupId:", req.SetupId)

	strsetupid := fmt.Sprintf("%d", req.SetupId)

	keys := fmt.Sprintf("setupid_%d", head.Uid)
	h.UserRedis.RedisSet(keys, strsetupid)

	common.DBUpdateSetupId(head.Uid, req.SetupId)
	return []byte(""), 0
}
