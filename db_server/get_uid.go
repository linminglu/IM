package db

import (
	"encoding/json"
	//	"fmt"
	//	"strconv"
	//	"strings"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type GetUidReq struct {
	PhoneNum string `json:"phonenum,omitempty"`
}

type GetUidResp struct {
	Uid uint64 `json:"uid,omitempty"`
}

func (h *DBHandler) GetUid(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	logger.Debug("GetUid jsonBody=", string(jsonBody))

	var req GetUidReq
	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	logger.Debug("GetUid req.PhoneNum=", req.PhoneNum)

	if req.PhoneNum == "" {
		logger.Error("phonenum is nil")
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	userInfo, err := common.GetUserInfoByPhoneNum(req.PhoneNum)
	if err != nil || userInfo == nil {
		logger.Error(err)
		return []byte(""), common.ERR_CODE_NO_USER
	}

	//	uid := fmt.Sprintf("%s",userInfo.Uid)
	//	uid := strconv.FormatUint(userInfo.Uid, 10)
	//	logger.Debug("uid=", uid)

	logger.Debug("userInfo.Uid=", userInfo.Uid)

	resp := GetUidResp{
		Uid: userInfo.Uid,
	}

	respBuf, err := json.Marshal(resp)

	return respBuf, 0

	//	if req.Appkey == "" || req.Cidlist == nil {
	//		logger.Error("nil arguments. Appkey: ", req.Appkey)
	//		return []byte(""), common.ERROR_CLIENT_BUG
	//	}
	//
	//	uidNum := len(req.Cidlist)
	//	logger.Info("request uid number: ", uidNum)
	//
	//	if uidNum <= 0 {
	//		logger.Error("request %d uids", uidNum)
	//		return []byte(""), common.ERROR_CLIENT_BUG
	//	}
	//
	//	//Get uids from redis
	//	uid_list := make([]string, uidNum)
	//	for i, v := range req.Cidlist {
	//		uid_list[i] = req.Appkey + "_" + v
	//	}

	//	varr, result := h.UserRedis.RedisMGet(uid_list)

	//	if result != 0 || varr[0] == nil {
	//		return []byte(""), common.ERR_CODE_SYS
	//	}
	//
	//	respstr := `{"uidlist":[`
	//	for _, v := range varr {
	//		uid, err := strconv.Atoi(v.(string))
	//		if err != nil {
	//			logger.Error("Atoi", v, " is err")
	//			respstr += "0,"
	//		} else {
	//			respstr += fmt.Sprintf("%d,", uid)
	//		}
	//	}

	//	respstr = strings.TrimRight(respstr, ",")
	//	respstr += "]}"

	//	return []byte(respstr), 0

}
