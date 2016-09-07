package kefu

import (
	"fmt"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

//import "strconv"

func (h *Handler) Current(ctx *web.Context, val string) {
	logger.Debug("Current start")

	retStr := ""

	defer func() {
		logger.Debug("return:", retStr)
		ctx.Write([]byte(retStr))
	}()

	logger.Debug("head:", ctx.Request.Header)
	logger.Debug("Form:", ctx.Request.Form)
	logger.Debug("PostForm:", ctx.Request.PostForm)

	retCookie, err := ctx.Request.Cookie("JSESSIONID")

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s","data":{}}`, err.Error())
		return
	}

	logger.Info("cookie:", retCookie.Value, err)

	app_key, errcode, _, account := h.Check(retCookie.Value)
	logger.Info(h.Check(retCookie.Value))

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

	logo, name := common.DBCSGetLogoName(app_key)

	retStr = fmt.Sprintf(`{"code":0,"err_msg":"", "data":{"account":"%s", "type":1, "app_key":"%s", "system_name":"%s","system_logo":"%s" }}`, account, app_key, name, logo)

	return
}
