package kefu

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type EnableReq struct {
	AppKey  string `json:"app_key,omitempty"`
	Account string `json:"account,omitempty"`
	Enable  int    `json:"enable,omitempty"`
}

func (h *Handler) Enable(ctx *web.Context, val string) {
	logger.Debug("Enable start")

	retStr := ""

	defer func() {
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

	var req EnableReq
	err = json.Unmarshal([]byte(jsonStr), &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, err.Error())
		return
	}

	if req.AppKey != app_key || len(req.Account) < 2 || req.Enable < 0 || req.Enable > 1 {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, err.Error())
		return
	}

	key := req.AppKey + "_" + req.Account
	v, result := h.UserRedis.RedisGet(key)

	if result != 0 {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "account not exist 1")
		return
	}

	if v == "" {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "account not exist 2")
		return
	}

	uid, err := strconv.ParseUint(v, 10, 64)

	if err != nil || uid == 0 {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "account not exist 3")
		return
	}

	userInfo := &common.UserInfo{uid, "", "", req.AppKey, req.Account, "", "", "", 0, 0, 0}
	nret := userInfo.DBUpdateCSEnable(req.Enable)

	if nret == 0 {
		retStr = fmt.Sprintf(`{"code":0,"err_msg":"%s"}`, "")
	} else {
		retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "system busy")
	}

	return
}
