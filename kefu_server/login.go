package kefu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type LoginReq struct {
	AppKey   string `json:"app_key,omitempty"`
	Password string `json:"password,omitempty"`
	Account  string `json:"account,omitempty"`
}

func (h *Handler) Login(ctx *web.Context, val string) {
	logger.Debug("Login start")

	retStr := ""

	defer func() {
		logger.Info("return :", retStr)
		ctx.Write([]byte(retStr))
	}()

	logger.Debug("head:", ctx.Request.Header)
	logger.Debug("Form:", ctx.Request.Form)
	logger.Debug("PostForm:", ctx.Request.PostForm)
	reqBuf := make([]byte, 1024)
	i := 0
	strLen := 0
	var err error

	for i < 4 {
		i++
		strLen, err = ctx.Request.Body.Read(reqBuf)

		if strLen > 0 {
			break
		}

		if err != nil {
			logger.Info("readlen :", strLen, "read req:", err)
			time.Sleep(time.Second)
			continue
		}
	}

	logger.Debug("req:", string(reqBuf[0:strLen]))

	var req LoginReq
	err = json.Unmarshal(reqBuf[0:strLen], &req)
	if err != nil {
		retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "json body error")
		return
	}

	userInfo := &common.UserInfo{0, "", "", req.AppKey, "", "", "", "", 0, 0, 0}

	uid, passwd := userInfo.DBCSLogin(req.AppKey, req.Account)
	logger.Info("DBCSLogin return uid %d, passwd %s", uid, passwd)

	if uid == 0 {
		retStr = fmt.Sprintf(`{"code":3001,"err_msg":"%s"}`, "invalid user")
		return
	}

	if passwd != req.Password {
		logger.Info("password error")
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "account or passwd error")
		return
	}

	//set cookie
	retCookie, err := ctx.Request.Cookie("JSESSIONID")

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s","data":{}}`, err.Error())
		return
	}

	cookiestr := fmt.Sprintf("%s|%d|%d|%s", req.AppKey, time.Now().UnixNano(), uid, req.Account)
	h.Session[retCookie.Value] = cookiestr

	ret := h.UserRedis.RedisSetEx(retCookie.Value, time.Second*86400, cookiestr)

	if ret != 0 {
		retStr = fmt.Sprintf(`{"code":3001,"err_msg":"%s"}`, "invalid user")
		return
	}

	logger.Info("login success")
	retStr = fmt.Sprintf(`{"code":0,"account":"%s", "type":1}`, req.Account)

	return
}
