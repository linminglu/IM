package kefu

import (
	"encoding/json"
	"fmt"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type HistoryReq struct {
	Itime   int    `json:"itime,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Account string `json:"account,omitempty"`
}

func (h *Handler) History(ctx *web.Context, val string) {
	logger.Debug("History start")

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

	var req HistoryReq
	err = json.Unmarshal(reqBuf[0:strLen], &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, "json body error")
		return
	}

	logger.Debug("req :", req)
	if req.Limit <= 0 {
		req.Limit = 300
	}

	if req.Itime < 1400000000 {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"itime %d error"}`, req.Itime)
		return
	}

	msgs := common.DBCSQueryMsgByTime(appkey, account, req.Itime, req.Itime+req.Limit)
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
