package rest_server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

func (h *Handler) Users(res http.ResponseWriter, req *http.Request) {
	maxCount := 100

	restResp := RestResp{State: 0, Msg: "ok"}

	defer func() {
		result, err := json.Marshal(restResp)
		if err != nil {
			res.Write([]byte(`{"state":500,"msg":"server err"}`))
			logger.Error(err)
		} else {
			res.Write(result)
		}
	}()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		restResp.State = 400
		restResp.Msg = "request json error"
		logger.Error(err)
		return
	}

	logger.Debug("rest api users req:", string(body))

	js, err := simplejson.NewJson(body)
	if err != nil {
		restResp.State = 400
		restResp.Msg = "request json error"
		logger.Error(err)
		return
	}

	arr, err := js.Array()
	if err != nil || len(arr) > maxCount {
		restResp.State = 400
		restResp.Msg = "request json error or json array too long"
		logger.Error(err)
		return
	}

	strSql := "insert into t_user_info (id, uid, reg_date, update_date, password) values "
	for i := 0; i < len(arr); i++ {
		uid, err := js.GetIndex(i).Get("uid").Int()

		if err != nil {
			restResp.State = 400
			restResp.Msg = "request json error"
			logger.Error(err)
			return
		}

		password, err := js.GetIndex(i).Get("password").String()

		if err != nil {
			restResp.State = 400
			restResp.Msg = "request json error"
			logger.Error(err)
			return
		}

		s := fmt.Sprintf("(0 , %d, now(), now(), '%s'), ", uid, password)
		strSql += s
	}

	strSql = strSql[:len(strSql) - 2]

	//	logger.Debug("strSql=", strSql)

	err = common.ProcExec(strSql)

	if err != nil {
		restResp.State = 401
		restResp.Msg = "add user fail, maybe user exists"
		logger.Error(err)
		return
	}
}
