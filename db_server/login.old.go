package db
//
//import (
//	"encoding/json"
//	//	"strconv"
//
//	"github.com/donnie4w/go-logger/logger"
//
//	"sirendaou.com/duserver/common"
//	"strconv"
//)
//
//type LoginReq struct {
//	Platform string `json:"platform,omitempty"`
//	Uid      int64  `json:"uid,omitempty"`
//	PhoneNum string `json:"phonenum,omitempty"`
//	Password string `json:"password,omitempty"`
//	Auto     uint64 `json:"auto,omitempty"`
//	SetupId  uint64 `json:"setupid,omitempty"`
//}
//
//type LoginResp struct {
//	Uid int64 `json:"uid,omitempty"`
//	Sid int   `json:"sid,omitempty"`
//}
//
//func (h *DBHandler) Login(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
//	logger.Debug("Login")
//	var req LoginReq
//
//	err := json.Unmarshal(jsonBody, &req)
//	if err != nil {
//		logger.Error("Unmarshal error:", err)
//		return []byte(""), common.ERROR_CLIENT_BUG
//	}
//
//	if req.Platform == "" || req.PhoneNum == "" || req.Password == "" {
//		logger.Error("nil arguments. Platform: ", req.Platform, " Phonenum: ", req.PhoneNum, " Password: ", req.Password)
//		return []byte(""), common.ERROR_CLIENT_BUG
//	}
//
//	if req.Platform == "a" || req.Platform == "A" {
//		tail.MsgId = uint64('a')
//	} else if req.Platform == "i" || req.Platform == "I" {
//		tail.MsgId = uint64('i')
//	} else if req.Platform == "w" || req.Platform == "W" {
//		tail.MsgId = uint64('w')
//	}
//
//	//	phonenum_uid_key := "du_phonenum_uid@" + req.PhoneNum
//	//	uidStr, code := h.UserRedis.RedisGet(phonenum_uid_key)
//	//	if code != 0 || uidStr == "" {
//	//		logger.Debug("user not exists Phonenum=", req.PhoneNum)
//	//		return []byte(""), common.ERR_CODE_NO_USER
//	//	}
//
//	//	tail.FromUid = userInfo.Uid
//
//	//
//	//	passKey := "Pwd_" + uidStr
//	//	pwd, _ := h.UserRedis.RedisGet(passKey)
//
//	userInfo, err := common.GetUserInfoByPhoneNum(req.PhoneNum)
//	if err != nil || userInfo == nil {
//		logger.Error(err)
//		return []byte(""), common.ERR_CODE_NO_USER
//	}
//
//	logger.Info("PhoneNum ", req.PhoneNum, " input passwd", req.Password, " right passwd:", userInfo.Pwd)
//
//	//	if uid&0xf == common.PT_KF {
//	//		logger.Info("uid ", uid, " is kefu")
//	//		return []byte(""), common.ERR_CODE_CS_LOGIN
//	//	}
//
//	if userInfo.Pwd != req.Password {
//		logger.Debug("PhoneNum ", req.PhoneNum, " error password")
//		return []byte(""), common.ERR_CODE_PASSWD
//	}
//
//	tail.FromUid = userInfo.Uid
//
//	//	if req.Auto == 1 && req.SetupId > 0 {
//	//		strSetupId, _ := h.UserRedis.RedisGet("setupid_" + uidStr)
//	//
//	//		logger.Info("uid ", uid, " input setupid", req.SetupId, " right cid:", strSetupId)
//	//
//	//		nSetupId, _ := strconv.ParseUint(strSetupId, 10, 64)
//	//		if nSetupId != req.SetupId {
//	//			return []byte(""), common.ERR_AUTOLOGIN_CONFLICT
//	//		}
//	//	}
//
//	uidStr := strconv.FormatUint(userInfo.Uid, 10)
//	passKey := "Pwd_" + uidStr
//	h.UserRedis.RedisSet(passKey,req.Password)
//
//	resp := LoginResp{
//		Uid: int64(userInfo.Uid),
//		Sid: int(tail.Sid),
//	}
//
//	respBuf, err := json.Marshal(resp)
//
//	return respBuf, 0
//}
