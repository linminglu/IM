package db

import (
	"encoding/json"
	//	"fmt"
	//	"strconv"
	"io/ioutil"
	"os"
	"time"
	//	"encoding/binary"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
	"strconv"
)

type RegReq struct {
	Platform string `json:"platform,omitempty"`
	DeviceId string `json:"did,omitempty"`
	PhoneNum string `json:"phonenum,omitempty"`
	Password string `json:"password,omitempty"`
	Login    int    `json:"login,omitempty"`
}

type RegResp struct {
	Uid int64 `json:"uid,omitempty"`
	Sid int   `json:"sid,omitempty"`
}

func (h *DBHandler) Register(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req RegReq
	err := json.Unmarshal(jsonBody, &req)

	isNewUser := true

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if req.PhoneNum == "" || req.Platform == "" || req.Password == "" {
		logger.Error("empty arguments. Platform: ", req.Platform, " Phonenum: ", req.PhoneNum, " Password: ", req.Password)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	userInfo, err := common.GetUserInfoByPhoneNum(req.PhoneNum)
	if userInfo != nil {
		logger.Info("phonenum has been registed phonenum=", req.PhoneNum)
		return []byte(""), common.ERR_CODE_PHONENUM_USED
	}

	//	// TODO check isNewUser: redis & db
	//	phonenum_uid_key := "du_phonenum_uid@" + req.PhoneNum
	//	uidInRedis, _ := h.UserRedis.RedisGet(phonenum_uid_key)
	//	if uidInRedis != "" {
	//		logger.Info("phonenum has been registed phonenum=", req.PhoneNum)
	//		return []byte(""), common.ERR_CODE_PHONENUM_USED
	//	} else {
	//		isNewUser = true
	//	}

	logger.Debug("Register isNewUser=", isNewUser)

	//	respCid := ""
	//	if req.Custom_id == "" {
	//		respCid = fmt.Sprintf("%d@daou", time.Now().UnixNano())
	//		req.Custom_id = respCid
	//		logger.Info("input cid is empty, create one:", respCid)
	//	}

	pt := uint64(0)

	if req.Platform == "a" || req.Platform == "A" {
		tail.MsgId = uint64('a')
		pt = common.PT_ANDROID
	} else if req.Platform == "i" || req.Platform == "I" {
		tail.MsgId = uint64('i')
		pt = common.PT_IOS
	} else if req.Platform == "w" || req.Platform == "W" {
		tail.MsgId = uint64('w')
		pt = common.PT_WP
	}

	if pt < common.PT_IOS || pt > common.PT_WP {
		logger.Error("Platform: ", req.Platform, " error")
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	//	key := req.Custom_id
	//	key := req.PhoneNum

	uid, err := getUid()
	if err != nil {
		logger.Error("getUid error", err.Error())
		return []byte(""), common.ERR_CODE_SYS
	}

	//	logger.Info("redis get", key, " uid ", uid)
	//
	//	if err != nil {
	//		logger.Error("redis get uid", key, " err :", err)
	//		return []byte(""), common.ERR_CODE_SYS
	//	}

	if isNewUser {
		if uid > 100000000 {
			return []byte(""), common.ERR_CODE_SYS
		}

		//		strUid := strconv.FormatUint(uid, 10)
		//
		//		h.UserRedis.RedisSet(key, strUid)
		//		h.UserRedis.RedisSAdd("Uidlist", strUid)
		//
		//		passKey := "Pwd_" + strUid
		//		h.UserRedis.RedisSet(passKey, req.Password)

		//		vk := "Cid_" + strUid
		//		h.UserRedis.RedisSet(vk, req.Custom_id)

		//		if h.UserRedis.RedisSet(phonenum_uid_key, strUid) != 0 {
		//			return []byte(""), common.ERR_CODE_SYS
		//		}

		userInfo := &common.UserInfo{
			Uid:      uid,
			Pwd:      req.Password,
			PhoneNum: req.PhoneNum,
			Platform: req.Platform,
			Did:      req.DeviceId,
			BaseInfo: "",
			ExInfo:   "",
			RegDate:  uint64(time.Now().Unix()),
			BV:       0,
			V:        0,
		}

		if err := userInfo.DBInsertUser(); err != nil {
			logger.Error(err)
			return []byte(""), common.ERR_CODE_SYS
		}

		if err := userInfo.DBInsertPhoneNum(); err != nil {
			logger.Error(err)
			return []byte(""), common.ERR_CODE_SYS
		}

		//插入新位置信息
		locInfo := &common.LocationInfo{uid, []float64{0, 0}, uint32(time.Now().Unix())}
		locInfo.NewLocation()
	} else {
		//		passkey := "Passwd_" + v
		//		passwd, _ := h.UserRedis.RedisGet(passkey)
		//		if passwd != req.Password {
		//			logger.Info("uid:", v, "Wrong password: ", passwd)
		//			return []byte(""), common.ERR_CODE_PASSWD
		//		}
	}

	tail.FromUid = uid

	resp := RegResp{int64(uid), int(tail.Sid)}
	respBuf, err := json.Marshal(resp)

	return respBuf, 0
}

// 从文件中读取uid，读取一次后uid+=1
func getUid() (uid uint64, err error) {
	f, err := os.OpenFile("../data/uid.dat", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("error e=", err.Error())
		return 0, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Error("error e=", err.Error())
		return 0, err
	}

	if len(data) == 0 {
		logger.Error("data is empty")
		uid = 100000
	} else {
		uid, err = strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			logger.Error("error", err.Error())
			return 0, err
		}
	}

	str := strconv.FormatUint(uid+1, 10)
	logger.Debug("str=", str)

	if _, err := f.Seek(0, 0); err != nil {
		logger.Error("error e=", err.Error())
		return 0, err
	}

	if _, err := f.WriteString(str); err != nil {
		logger.Error("error e=", err.Error())
		return 0, err
	}

	return uid, err
}
