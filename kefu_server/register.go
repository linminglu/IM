package kefu

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type RegExinfo struct {
	Image_id  string `json:"image_id,omitempty"`
	Nick_name string `json:"nick_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Tel       string `json:"tel,omitempty"`
}

type RegReq struct {
	Account  string    `json:"account,omitempty"`
	Password string    `json:"password,omitempty"`
	Exinfo   RegExinfo `json:"ext_info,omitempty"`
}

func (h *Handler) Reg(appkey string, req RegReq) int {
	key := appkey + "_" + req.Account
	v, result := h.UserRedis.RedisGet(key)

	if result != 0 {
		return 2
	}

	isNew := false
	if v == "" {
		v = h.UserRedis.RedisRPop(common.REDIS_UID_POOL)

		if v == "" {
			logger.Error("redis get", common.REDIS_UID_POOL, " for free uid from redis err.")
			return 2
		} else {
			logger.Info("RedisRPop:", common.REDIS_UID_POOL, v)
			isNew = true
		}
	}

	logger.Info("redis get", common.REDIS_UID_POOL, " uid ", v)

	uid, _ := strconv.ParseUint(v, 10, 64)

	if uid == 0 {
		return 2
	}

	if !isNew && (uid&0xf != 6) {
		logger.Info("account ", req.Account, uid, " has been used")
		return 1
	}

	longUid := uint64(0)
	if isNew {
		longUid = common.GetLongUid(uid, uint64(common.PT_KF))
	} else {
		longUid = uid
	}

	strUid := fmt.Sprintf("%d", longUid)

	h.UserRedis.RedisSet(key, strUid)

	baseInfo := fmt.Sprintf(`{"nick_name":"%d", "image_id":"%s","email":"%d","%d","tel":"%d"}`, req.Exinfo.Nick_name, req.Exinfo.Image_id, req.Exinfo.Email, req.Exinfo.Tel)

	userInfo := &common.UserInfo{longUid, req.Password, appkey, req.Account, "kf", "kefu", baseInfo, "", uint64(time.Now().Unix()), 0, 0}

	if err := userInfo.DBInsertUser(); err != nil {
		logger.Error(err)
	}

	//func (user *UserInfo) DBInsertCS(image, tel, name, email string) int {
	userInfo.DBInsertCS(req.Exinfo.Image_id, req.Exinfo.Tel, req.Exinfo.Nick_name, req.Exinfo.Email)

	return 0
}

func (h *Handler) Register(ctx *web.Context, val string) {
	logger.Debug("Init start")

	retStr := ""

	defer func() {
		logger.Info("return:", retStr)
		ctx.Write([]byte(retStr))
	}()

	app_key, master_key, err := ctx.GetBasicAuth()

	if err != nil {
		retStr = `{"code":1001,"err_msg":""}`
		return
	}

	logger.Info("appkey:", app_key, "master_key:", master_key)

	ret := h.CheckAuth(app_key, master_key)

	if ret != 0 {
		retStr = `{"code":1001,"err_msg":""}`
		return
	}

	reqBuf := make([]byte, 1024)
	strLen, _ := ctx.Request.Body.Read(reqBuf)

	jsonStr := string(reqBuf[0:strLen])
	logger.Debug("req:", jsonStr)

	var req RegReq
	err = json.Unmarshal([]byte(jsonStr), &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, err.Error())
		return
	}

	if len(req.Account) < 2 || len(req.Password) < 2 || len(req.Exinfo.Nick_name) < 1 {
		retStr = fmt.Sprintf(`{"code":1003, "err_msg":"%s"}`, "")
		return
	}

	ret = h.Reg(app_key, req)
	if ret != 0 {
		switch ret {
		case 1:
			retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "appkey is error")
		default:
			retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "")
			return
		}
	}

	retStr = fmt.Sprintf(`{"code":0,"err_msg":"","data":{"account":"%s"}}`, req.Account)
	return
}
