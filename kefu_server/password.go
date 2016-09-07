package kefu

import (
	"encoding/json"
	"fmt"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

//import "strconv"

type PwdReq struct {
	Password string `json:"password,omitempty"`
}

func (h *Handler) Password(ctx *web.Context, val string) {

	logger.Debug("Password start")

	retStr := ""

	defer func() {
		logger.Debug("return:", retStr)
		ctx.Write([]byte(retStr))
	}()

	reqBuf := make([]byte, 1024)
	strLen, err := ctx.Request.Body.Read(reqBuf)

	logger.Debug("req:", string(reqBuf[:]))

	var req PwdReq
	err = json.Unmarshal(reqBuf[0:strLen], &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "json body error")
		return
	}

	retCookie, err := ctx.Request.Cookie("JSESSIONID")

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s","data":{}}`, err.Error())
		return
	}

	logger.Info("cookie:", retCookie.Value, err)

	app_key, errcode, uid, account := h.Check(retCookie.Value)
	logger.Info(app_key, errcode, uid, account)

	if errcode != 0 {
		switch errcode {
		case 1:
			retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s","data":{}}`, "")
		case 2:
			retStr = fmt.Sprintf(`{"code":1002,"err_msg":"%s","data":{}}`, "")
		default:
			retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s","data":{}}`, "")

		}
		return
	}

	ret := common.DBCSResetPwd(uid, req.Password)
	if ret != 0 {
		retStr = fmt.Sprintf(`{"code":1000,"err_msg":""}`)
		return

	}

	retStr = fmt.Sprintf(`{"code":0,"err_msg":""}`)

	return
}
