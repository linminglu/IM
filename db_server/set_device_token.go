package db

import (
	"encoding/json"
	"fmt"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type DeviceTokenReq struct {
	DeviceToken string `json:"devicetoken,omitempty"`
}

func (h *DBHandler) SetDeviceToken(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req DeviceTokenReq

	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	common.DBUniInsertToken(head.Uid, req.DeviceToken)

	if req.DeviceToken != "" {
		common.DBClearExtraToken(head.Uid, req.DeviceToken)
	}

	key := fmt.Sprintf("token_%d", head.Uid)
	if len(req.DeviceToken) > 32 {
		h.UserRedis.RedisSet(key, req.DeviceToken)
	} else {
		h.UserRedis.RedisDel(key)
	}

	return []byte(""), 0
}
