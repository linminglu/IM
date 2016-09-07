package file_server

import (
	//	"crypto/hmac"
	//	"crypto/sha1"
	//	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"
	"github.com/qiniu/api/rs"
)

var page2 = `<html><head></head><body>
<form method="POST" enctype="multipart/form-data" action="/file/upload">
fromuid:<input type="text" name="fromuid" /> <br/>
hash:<input type="text" name="hash" /> <br/>
vercode:<input type="text" name="vercode"/>
<br/>
<input type="submit" />
</form>
</body></html>`

var page3 = `<html><head></head><body>
<form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
  <input name="key" type="text" value="<Your file name in qiniu>">
  <input name="token" type="text" value="<Your uptoken from server>">
  <input name="file" type="file"/>
<input type="submit" />
</form>
/<body></html>`

func uploadHandlerGet(ctx *web.Context, val string) {
	logger.Info("Get upload")

	_, err := ctx.Write([]byte(page2))
	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}

func uploadHandlerGet2(ctx *web.Context, val string) {
	logger.Info("Get upload")

	_, err := ctx.Write([]byte(page3))
	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}

//func GetToken(md5, fromUid string) string {
//	callbackUrl := *g_CallbackUrl
////	SecretKey := "pFUqIQPOTgatStZEaXdwTTZVwFkIiDzJWzmrSAyT"
////	AccessKey := "xs5omPdRjYfP2T3OelgoPBEl8aoEMKc4VUQ1cJ0P"
//
//	deadline := time.Now().Unix() + 180
//
//	keys := fmt.Sprintf("%s%s%s", string([]byte(md5)[0:5]), fromUid, string([]byte(md5)[5:]))
//	scope := "du:" + keys
//
//	//"callbackUrl":"%s".
//	//"callbackBody ":"hash=$(etag)&uid=%s&fid=%s",
//	putPolicy := fmt.Sprintf(`{"scope":"%s","deadline":%d,"returnBody":"{\"hash\":$(etag)}"}`,
//		scope, deadline, callbackUrl, fromUid, keys)
//
//	encodedPutPolicy := base64.StdEncoding.EncodeToString([]byte(putPolicy))
//	logger.Debug("putPolicy", putPolicy, "encodedPutPolicy:", encodedPutPolicy)
//	h := hmac.New(sha1.New, []byte(SecretKey))
//	h.Write([]byte(encodedPutPolicy))
//	encodedSign := base64.StdEncoding.EncodeToString(h.Sum(nil))
//
//	uploadToken := AccessKey + ":" + encodedSign + ":" + encodedPutPolicy
//	logger.Debug("token:", uploadToken)
//	return uploadToken
//}

func GetToken2(md5, fromUid string) (string, string) {
	uid, _ := strconv.Atoi(fromUid)
	key := fmt.Sprintf("%s%s%sfe%d", string([]byte(md5)[0:5]), fromUid, string([]byte(md5)[5:]), uid%99973)
	logger.Info("GetToken2", fromUid, key)
	scope := "duserver:" + key
	callbackBody := fmt.Sprintf("hash=$(etag)&uid=%s&fid=%s", fromUid, key)
	expires := time.Now().Unix() + 180
	putPolicy := rs.PutPolicy{
		Scope:        scope,
		CallbackUrl:  *g_CallbackUrl,
		CallbackBody: callbackBody,
		//ReturnUrl:   returnUrl,
		//ReturnBody:  returnBody,
		//AsyncOps:    asyncOps,
		//EndUser:     endUser,
		Expires: uint32(expires),
	}
	return putPolicy.Token(nil), key
}

func UploadHandler(ctx *web.Context, val string) {
	ctx.Request.ParseMultipartForm(1024)

	fromUid := ctx.Request.FormValue("fromuid")
	md5 := ctx.Request.FormValue("hash")
	verCode := ctx.Request.FormValue("vercode")
	logger.Info("fromuid:", fromUid, "vercode:", verCode, "hash:", md5)
	retStr := ""
	for {
		if len(verCode) != 32 {
			ctx.Unauthorized()
			retStr = `{"errcode":1000}`
			break
		}
		token, key := GetToken2(md5, fromUid)
		retStr = fmt.Sprintf(`{"errcode":0,"token":"%s","url":"%s", "key":"%s"}`, token, *g_UploadUrl, key)
		break
	}

	_, err := ctx.Write([]byte(retStr))
	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}

func CallBack(ctx *web.Context, val string) {
	ctx.Request.ParseMultipartForm(1024)

	md5 := ctx.Request.FormValue("hash")
	uid := ctx.Request.FormValue("uid")
	fid := ctx.Request.FormValue("fid")

	logger.Info("--------------CallBack:", uid, md5, fid)

	retStr := fmt.Sprintf(`{"fid":"%s"}`, fid)

	_, err := ctx.Write([]byte(retStr))

	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}
