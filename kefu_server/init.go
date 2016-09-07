package kefu

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"

	"sirendaou.com/duserver/common"
)

type InitReq struct {
	Logo string `json:"logo,omitempty"`
	Name string `json:"name,omitempty"`
}

const (
	msg = `
0	成功
1000	其它未知异常
1001	鉴权失败
1002	登录超时
1003	参数异常
1004	权限不足
1005	服务端超时
2001	表示由于'CustomerService'这个ID存在导致初始化失败
3001	无效帐号，可能是不存在，也可能是禁用之类
3002	密码有误
4001	错误的目标用户
4002	错误的消息类型
4003	消息体超长`
)

func (h *Handler) Check(cookie string) (string, int, uint64, string) {
	v, ok := h.Session[cookie]
	if !ok {
		v, _ := h.UserRedis.RedisGet(cookie)

		if len(v) <= 1 {
			logger.Info("redis get cookie ", cookie, "fail")
			return "", 1, 0, ""
		} else {
			h.Session[cookie] = v
		}
	}

	cookie = v
	cookies := strings.Split(cookie, "|")

	if len(cookies) != 4 {
		return "", 1, 0, ""
	}

	app_key, sTime, sUid, account := cookies[0], cookies[1], cookies[2], cookies[3]

	ntime, _ := strconv.ParseInt(sTime, 10, 64)
	uid, _ := strconv.ParseUint(sUid, 10, 64)

	logger.Info("appkey:", app_key, "lasttime", ntime, "uid", uid)

	//if ntimee+1800000000000 < time.Now().UnixNano() {
	//	logger.Info("cookie time out")
	//	return "", 2, 0, ""
	//}

	return app_key, 0, uid, account
}

func (h *Handler) CheckAuth(appkey, master_sec string) int {
	return 0
}

func (h *Handler) Init(ctx *web.Context, val string) {

	logger.Debug("Init start")

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

	var req InitReq
	err = json.Unmarshal([]byte(jsonStr), &req)

	if err != nil {
		retStr = `{"code":1001,"err_msg":""}`
		return
	}

	common.DBCSUpdateLogoName(app_key, req.Logo, req.Name)
	retStr = `{"code":0,"err_msg":""}`
	return
}
