package kefu

import (
	"encoding/json"
	"fmt"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type MsgListReq struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Account  string `json:"account,omitempty"`
}

func (h *Handler) MsgList(ctx *web.Context, val string) {
	logger.Debug("Init start")

	retStr := ""

	defer func() {
		ctx.Write([]byte(retStr))
	}()

	retCookie, ok := ctx.Request.Cookie("JSESSIONID")
	logger.Info("cookie:", retCookie.Value, ok)

	appkey, errcode, _, account := h.Check(retCookie.Value)

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

	reqBuf := make([]byte, 1024)
	strLen, err := ctx.Request.Body.Read(reqBuf)

	logger.Debug("req:", string(reqBuf[:strLen]))

	var req MsgListReq
	err = json.Unmarshal(reqBuf[0:strLen], &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "json body error")
		return
	}

	if req.Page <= 0 {
		req.Page = 0
	} else {
		req.Page--
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	msgs := common.DBCSQueryMsg(appkey, account, req.Page, req.PageSize)
	retStr = fmt.Sprintf(`{"code":0,"err_msg":"","data":{"total":%d, "list":[`, len(msgs))
	for i, msg := range msgs {
		if i == 0 {
			retStr += msg
		} else {
			retStr = retStr + "," + msg
		}
	}

	retStr += `] }}`

	return
}
