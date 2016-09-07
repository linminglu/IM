package db

import (
	"encoding/json"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
	"strconv"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"time"
)

type LoginReq struct {
	Platform string `json:"platform,omitempty"`
	Uid      uint64 `json:"uid,omitempty"`
	Password string `json:"password,omitempty"`
	SetupID  string `json:"setupid,omitempty"`
	DeviceID string `json:"deviceid,omitempty"`
}

type LoginResp struct {
	Uid int64 `json:"uid,omitempty"`
	Sid int   `json:"sid,omitempty"`
	Token string `json:"token,omitempty"`
}

type TokenLoginReq struct {
	Token    string `json:"token,omitempty"`
}

func (h *DBHandler) Login(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req LoginReq

	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if req.Platform == "" || req.Password == "" {
		logger.Error("parames error. Platform: ", req.Platform, " Password: ", req.Password)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	req.Platform = strings.ToLower(req.Platform)
	req.Platform = req.Platform[0:1]

	if req.Platform == "a" || req.Platform == "A" {
		tail.MsgId = uint64('a')
	} else if req.Platform == "i" || req.Platform == "I" {
		tail.MsgId = uint64('i')
	} else if req.Platform == "w" || req.Platform == "W" {
		tail.MsgId = uint64('w')
	} else if strings.HasPrefix(req.Platform, "p") || strings.HasPrefix(req.Platform, "P") {
		tail.MsgId = uint64('p')
	}

	userInfo, err := common.GetUserInfoByUid(req.Uid)
	if err != nil || userInfo == nil {
		logger.Error("login error:", err)
		return []byte(""), common.ERR_CODE_NO_USER
	}

	logger.Info("input passwd", req.Password, " right passwd:", userInfo.Pwd)
	if userInfo.Pwd != req.Password {
		return []byte(""), common.ERR_CODE_PASSWD
	}

	tail.FromUid = userInfo.Uid

	uidStr := strconv.FormatUint(userInfo.Uid, 10)
	passKey := "Pwd_" + uidStr
	h.UserRedis.RedisSet(passKey, req.Password)

	resp := LoginResp{
		Uid: int64(userInfo.Uid),
		Sid: int(tail.Sid),
		Token:common.DBCreateMobileLoginToken(req.Platform, userInfo.Uid, req.Password, "", "", ),
	}

	respBuf, err := json.Marshal(resp)
	if err != nil {
		logger.Error("login error:", err)
		return []byte(""), common.ERR_CODE_SYS
	}

	return respBuf, 0
}

func (h *DBHandler) LoginWithToken(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req TokenLoginReq

	err := json.Unmarshal(jsonBody, &req)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	// 密钥
	var keyText = "astaxie12798akljzmknm.ahkjkljl;k"
	var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	c, err := aes.NewCipher([]byte(keyText))
	if err != nil {
		logger.Error("NewCipher error:", err)
		return []byte(""), common.ERROR_UNKNOWN
	}

	// 解密
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	jsonString := make([]byte, len(req.Token))
	cfbdec.XORKeyStream(jsonString, []byte(req.Token))

	var tokenLoginBody common.TokenLoginBody

	err = json.Unmarshal(jsonString, &tokenLoginBody)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	// 参数检查
	if tokenLoginBody.Uid == 0 {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	// 参数检查
	if tokenLoginBody.PlatformType == nil || len(tokenLoginBody.PlatformType) {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	// 参数检查
	if tokenLoginBody.CoreToken == nil || len(tokenLoginBody.CoreToken) {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	uidStr := strconv.FormatUint(tokenLoginBody.Uid, 10)
	tokenKey := "LoginToken_" + uidStr

	realCoreToken, result := h.UserRedis.RedisGet(tokenKey)
	if result != 0 || realCoreToken == nil || realCoreToken == "" {
		realCoreToken = common.DBGetMobileCoreTokenWithUID(tokenLoginBody.Uid)
	}

	// token匹配错误
	if realCoreToken != tokenLoginBody.CoreToken {
		return []byte(""), common.ERR_TOKENLOGIN_UNKNOWN
	}

	tokenTimeKey := "LoginToken_time_" + uidStr

	tokenTime, result := h.UserRedis.RedisGet(tokenTimeKey)
	if result != 0 || tokenTime == nil || tokenTime == "" {
		tokenTime = common.DBGetMobileTokenTimeWithUID(tokenLoginBody.Uid)
	}

	// 超过5天
	// 5 * 24 * 3600 * 1,000,000,000 =
	if tokenTime != 0 && (time.Now().UnixNano() - tokenTime) > 432000000000000 {
		return []byte(""), common.ERR_TOKENLOGIN_EXPIRED
	}

	userInfo, err := common.GetUserInfoByUid(tokenLoginBody.Uid)
	if err != nil || userInfo == nil {
		logger.Error("login error:", err)
		return []byte(""), common.ERR_CODE_NO_USER
	}

	passKey := "Pwd_" + uidStr
	passwd, result := h.UserRedis.RedisGet(passKey)
	if result != 0 || passwd == nil || passwd == "" {
		realCoreToken = common.DBGetMobileCoreTokenWithUID(tokenLoginBody.Uid)
	}

	resp := LoginResp{
		Uid: int64(userInfo.Uid),
		Sid: int(tail.Sid),
		Token:common.DBCreateMobileLoginToken(tokenLoginBody.PlatformType, userInfo.Uid, userInfo.Pwd, "", "", ),
	}

	respBuf, err := json.Marshal(resp)

	return respBuf, 0
}