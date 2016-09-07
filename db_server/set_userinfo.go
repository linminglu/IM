package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type UserInfoSet struct {
	Did      string `json:"did,omitempty"`
	Baseinfo string `json:"baseinfo,omitempty"`
	Exinfo   string `json:"exinfo,omitempty"`
}

type UserInfoSetResp struct {
	Bv uint64
	V  uint64
}

func (this *UserInfoSetResp) MarshalJSON() ([]byte, error) {

	if this.Bv > 0 && this.V > 0 {
		return json.Marshal(map[string]interface{}{
			"bv": this.Bv,
			"v":  this.V,
		})
	} else if this.Bv > 0 {
		return json.Marshal(map[string]interface{}{
			"bv": this.Bv,
		})
	} else if this.V > 0 {
		return json.Marshal(map[string]interface{}{
			"v": this.V,
		})
	}

	return []byte("{}"), nil
}

func (h *DBHandler) SetUserInfo(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {

	var req UserInfoSet
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	uidstr := strconv.FormatUint(head.Uid, 10)

	exV := uint64(0)
	bV := uint64(0)

	if req.Did != "" {
		result := h.UserRedis.RedisSetEx("Did_"+uidstr, time.Second*common.EXPIRE_TIME, req.Did)
		if result != 0 {
			logger.Info("redis setex ", "Did_"+uidstr, "fail")
			return []byte(""), common.ERR_CODE_SYS
		}
		v, _ := h.UserRedis.RedisGet("Did_" + uidstr)
		logger.Info("Set Key", "Did_", uidstr, " value ", v, " success.")
	}

	if req.Baseinfo != "" {

		result := h.UserRedis.RedisSetEx("Baseinfo_"+uidstr, time.Second*common.EXPIRE_TIME, req.Baseinfo)
		if result != 0 {
			logger.Info("redis setex ", "Baseinfo_"+uidstr, "fail")
			return []byte(""), common.ERR_CODE_SYS
		}
		v, _ := h.UserRedis.RedisGet("Baseinfo_" + uidstr)
		logger.Info("Set Key", "Baseinfo_", uidstr, " value ", v, " success.")

		bV = uint64(time.Now().Unix())
	}

	if req.Exinfo != "" {
		result := h.UserRedis.RedisSetEx("Exinfo_"+uidstr, time.Second*common.EXPIRE_TIME, req.Exinfo)
		if result != 0 {
			logger.Info("redis setex ", "Exinfo_"+uidstr, "fail")
			return []byte(""), common.ERR_CODE_SYS
		}
		v, _ := h.UserRedis.RedisGet("Exinfo_" + uidstr)
		logger.Info("Set Key", "Exinfo_", uidstr, " value ", v, " success.")

		exV = uint64(time.Now().Unix())

	}

	if bV > 0 {
		strtime := fmt.Sprintf("%d", bV)
		result := h.UserRedis.RedisSet("BV_"+uidstr, strtime)
		if result != 0 {
			logger.Info("redis setex ", "BV_"+uidstr, "fail")
		}
	}

	if exV > 0 {
		strtime := fmt.Sprintf("%d", exV)
		result := h.UserRedis.RedisSet("V_"+uidstr, strtime)
		if result != 0 {
			logger.Info("redis setex ", "V_"+uidstr, "fail")
		}
	}

	/*
		type UserInfo struct {
			Uid      uint64 `json:"uid"`
			Passwd   string `json:"passwd"`
			Appkey   string `json:"appkey"`
			Cid      string `json:"cid"`
			Platform string `json:"platform"`
			Did      string `json:"did"`
			Baseinfo string `json:"baseinfo"`
			Exinfo   string `json:"exinfo"`
			RegDate  uint64 `json:"regdate"`
			BV       uint64 `json:"bv"`
			V        uint64 `json:"v"`
		}*/

	userInfo := &common.UserInfo{Uid: head.Uid, Did: req.Did, BaseInfo: req.Baseinfo, ExInfo: req.Exinfo, BV: bV, V: exV}
	userInfo.DBUpdateUserInfo()

	resp := UserInfoSetResp{bV, exV}
	strResp, err := resp.MarshalJSON()

	if err != nil {
		return []byte(""), common.ERR_CODE_SYS
	} else {
		return strResp, 0
	}
}
