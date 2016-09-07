package kefu

import (
	"encoding/json"
	"fmt"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type ListReq struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

func (h *Handler) List(ctx *web.Context, val string) {
	logger.Debug("List start")

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

	var req ListReq
	err = json.Unmarshal([]byte(jsonStr), &req)

	if err != nil {
		retStr = fmt.Sprintf(`{"code":1003,"err_msg":"%s"}`, err.Error())
		return
	}

	if req.Page > 0 {
		req.Page--
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	userInfo := &common.UserInfo{}
	ret, infos := userInfo.DBGetCSList(req.Page, req.PageSize)

	if ret != 0 {
		retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "system busy")
		return
	}

	total := len(infos)
	retStr = fmt.Sprintf(`{"code":0,"err_msg":"","data":{"total":%d, "list":[`, total)
	/*
	   type CSInfo struct {
	   	Uid     uint64
	   	Passwd  string
	   	Appkey  string
	   	Account string
	   	Name    string
	   	Image   string
	   	Tel     string
	   	Email   string
	   	Enable  int
	   }
	*/
	for i, info := range infos {
		tempstr := fmt.Sprintf(`{"cid":"%s","nick_name":"%s", "image_id":"%s", "email":"%s", "tel":"%s","enable":%d}`,
			info.Account, info.Name, info.Image, info.Email, info.Tel, info.Enable)
		retStr += tempstr
		if i < total-1 {
			retStr += ","
		}
	}

	retStr += "]}}"

	return
}
