package db

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type UserInfoReq struct {
	UidList      []uint64 `json:"uidlist,omitempty"`
	PropertyList []string `json:"propertylist,omitempty"`
}

func (h *DBHandler) GetUserInfo(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req UserInfoReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	uidNum := len(req.UidList)
	propertyNum := len(req.PropertyList)

	logger.Info("request uid number: ", uidNum, " property number:", propertyNum)
	logger.Info("req.PropertyList: ", req.PropertyList)

	if uidNum <= 0 || propertyNum <= 0 {
		logger.Error("request %d uids %d properties", uidNum, propertyNum)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	req.PropertyList = append(req.PropertyList, "uid")

	infoMap := common.DBGetUserInfo(req.UidList, req.PropertyList)

	for k, v := range infoMap {
		logger.Info("index: ", k, " Info: ", v)
	}

	resp := "{\"infolist\":["
	for _, v := range infoMap {
		resp += fmt.Sprintf("{%s},", v)
	}
	resp = strings.TrimRight(resp, ",")
	resp += "]}"

	return []byte(resp), 0
}
