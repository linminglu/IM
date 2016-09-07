package db
//
//
//import (
//	"encoding/json"
//	"fmt"
//	"strconv"
//	"time"
//
//	"github.com/donnie4w/go-logger/logger"
//
//	"sirendaou.com/duserver/common"
//)
//
//func (h *DBHandler) PureReg(head common.PkgHead, jsonbody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
//	var req RegReq
//	err := json.Unmarshal(jsonbody, &req)
//
//	isNewUser := false
//
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
//	pt := uint64(0)
//
//	if req.Platform == "a" || req.Platform == "A" {
//		tail.MsgId = uint64('a')
//		pt = common.PT_ANDROID
//	} else if req.Platform == "i" || req.Platform == "I" {
//		tail.MsgId = uint64('i')
//		pt = common.PT_IOS
//	} else if req.Platform == "w" || req.Platform == "W" {
//		tail.MsgId = uint64('w')
//		pt = common.PT_WP
//	}
//
//	if pt < common.PT_IOS || pt > common.PT_WP {
//		logger.Error("Platform: ", req.Platform, " error")
//		return []byte(""), common.ERROR_CLIENT_BUG
//	}
//
////	//检查Appkey
////	appkey := "appinfo_" + req.Appkey
////	v, result := h.AppRedis.RedisHGet(appkey, "dev_id")
////	if result != 0 {
////		logger.Error("Unable to find appkey ", req.Appkey)
////		return []byte(""), common.ERROR_OUT_OF_REACH
////	}
////	developerID, err := strconv.Atoi(v)
////	if err != nil {
////		logger.Error("redis get developer ID ", v, " err :", err)
////		return []byte(""), common.ERR_CODE_SYS
////	}
////
////	key := req.Appkey + "_" + req.Custom_id
////
////	v, result = h.UserRedis.RedisGet(key)
//
//	if result != 0 {
//		return []byte(""), common.ERR_CODE_SYS
//	}
//
//	if v == "" {
//		if req.Login > 0 {
//			return []byte(""), common.ERR_CODE_CID_EXIST
//		}
//
//		v = h.UserRedis.RedisRPop(common.REDIS_UID_POOL)
//
//		if v == "" {
//			logger.Error("redis get", key, " for free uid from redis err.")
//			return []byte(""), common.ERR_CODE_SYS
//		} else {
//			logger.Info("RPOP", common.REDIS_UID_POOL, v)
//		}
//
//		isNewUser = true
//	} else {
//		if req.Login == 0 {
////			logger.Info(req.Custom_id, " is exist")
//			return []byte(""), common.ERROR_ACCOUNT
//		}
//	}
//
//	logger.Info("redis get", key, " uid ", v)
//
//	uid, err := strconv.ParseUint(v, 10, 64)
//	longuid := uid
//
//	if err != nil {
//		logger.Error("redis get uid", key, " err :", err)
//		return []byte(""), common.ERR_CODE_SYS
//	}
//
//	if isNewUser {
//		//check appkey
//
//		if uid > 100000000 {
//			return []byte(""), common.ERR_CODE_SYS
//		}
//
////		appid := 0
////		appid, ok := h.AppkeyMap[req.Appkey]
////		if !ok {
////			appid = common.DBLoadAppIDBykey(req.Appkey)
////		}
////		if appid == 0 {
////			logger.Error("req.Appkey ", req.Appkey, " err cannot find")
////			return []byte(""), common.ERR_APPKEY
////		}
//
////		longuid = common.GetLongUid(uid,uint64(pt))
//		strUid := fmt.Sprintf("%d", uid)
//
////		h.UserRedis.RedisSet(key, strUid)
//
//		setKey :=  "du_uidlist"
//		h.UserRedis.RedisSAdd(setKey, strUid)
//
//		passkey := "Pwd_" + strUid
//		h.UserRedis.RedisSet(passkey, req.Password)
//
////		vk := "Cid_" + strUid
////		h.UserRedis.RedisSet(vk, req.Customd)
////
////		vk = "Appkey_" + strUid
////		h.UserRedis.RedisSet(vk, req.Appkey)
//
//		userInfo := &common.UserInfo{longuid, req.Password,req.Platform, req.DeviceId, "", "", uint64(time.Now().Unix()), 0, 0}
//		userInfo.DBInsertUser()
//
////		userinfo.SyncRegToDataCenter(*exchange, "rk-reg")
//
//		//插入新位置信息
//		locInfo := &common.LocationInfo{longuid,[]float64{0, 0}, uint32(time.Now().Unix())}
//		locInfo.NewLocation()
//	} else {
//		passKey := "Pwd_" + strUid
//		passwd, _ := h.UserRedis.RedisGet(passKey)
//
//		uid, _ := strconv.Atoi(uid)
//
//		if uid&0xf == common.PT_KF {
//			logger.Info("uid:", v, "is kefu")
//			return []byte(""), common.ERR_CODE_CS_LOGIN
//		}
//
//		if passwd != req.Password {
//			logger.Info("uid:", v, "Wrong password: ", passwd)
//			return []byte(""), common.ERR_CODE_PASSWD
//		}
//	}
//
//	tail.FromUid = uint64(longuid)
//
//	resp := RegResp{int64(longuid), int(tail.Sid), ""}
//
//	respBuf, err := json.Marshal(resp)
//
//	return respBuf, 0
//}
