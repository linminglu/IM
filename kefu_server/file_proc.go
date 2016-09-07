package kefu

import (
	"fmt"
	"strconv"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"
	"github.com/qiniu/api/fop"
	"github.com/qiniu/api/rs"
)

const SecretKey = "pFUqIQPOTgatStZEaXdwTTZVwFkIiDzJWzmrSAyT"
const AccessKey = "xs5omPdRjYfP2T3OelgoPBEl8aoEMKc4VUQ1cJ0P"

func GetUrl(fid string, url string) string {
	//deadline := time.Now().Unix() + 180
	baseUrl := rs.MakeBaseUrl(url, fid)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}

func GetImageUrl(w, h int, fid, url string) string {
	baseurl := rs.MakeBaseUrl(url, fid)
	var view = fop.ImageView{2, w, h, 80, "jpg"}
	imageurl := view.MakeRequest(baseurl)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(imageurl, nil)
}

func (handle *Handler) ImageDownload(ctx *web.Context, val string) {
	retStr := ""
	w := 100
	h := 100

	/*	defer func() {
			logger.Info("return:", retstr)
			ctx.Write([]byte(retstr))
		}()
	*/
	retCookie, ok := ctx.Request.Cookie("JSESSIONID")
	logger.Info("cookie:", retCookie, ok)

	logger.Info(handle.Check(retCookie.Value))
	_, errcode, _, _ := handle.Check(retCookie.Value)
	logger.Debug("error:", errcode)

	if errcode != 0 {
		switch errcode {
		case 1:
			retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s"}`, "")
		case 2:
			retStr = fmt.Sprintf(`{"code":1002,"err_msg":"%s"}`, "")
		default:
			retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "")
		}
		return
	}

	ctx.Request.ParseMultipartForm(1024)

	fid := ctx.Request.FormValue("fileID")
	sw := ctx.Request.FormValue("width")
	sh := ctx.Request.FormValue("high")

	w, _ = strconv.Atoi(sw)
	h, _ = strconv.Atoi(sh)

	logger.Info("fid:", fid, "w:", w, "h:", h)

	if len(fid) < 32 {
		ctx.NotFound("not found")
		//		retstr = `{"errcode":1002}`
		return
	}

	url := GetImageUrl(w, h, fid, *g_DownloadUrl)

	retStr = fmt.Sprintf(`{"code":0,"data":{"url":"%s"}}`, url)

	logger.Info("return:", retStr)

	_, err := ctx.Write([]byte(retStr))

	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}

func (handle *Handler) FileDownload(ctx *web.Context, val string) {
	retStr := ""

	retCookie, ok := ctx.Request.Cookie("JSESSIONID")
	logger.Info("cookie:", retCookie, ok)

	logger.Info(handle.Check(retCookie.Value))
	_, errcode, _, _ := handle.Check(retCookie.Value)
	logger.Debug("error:", errcode)

	if errcode != 0 {
		switch errcode {
		case 1:
			retStr = fmt.Sprintf(`{"code":1001,"err_msg":"%s"}`, "")
		case 2:
			retStr = fmt.Sprintf(`{"code":1002,"err_msg":"%s"}`, "")
		default:
			retStr = fmt.Sprintf(`{"code":1000,"err_msg":"%s"}`, "")

		}
		return
	}

	ctx.Request.ParseMultipartForm(1024)

	fid := ctx.Request.FormValue("fileID")

	logger.Info("get file fid:", fid)

	if len(fid) < 32 {
		ctx.NotFound("not found")
		return
	}

	url := GetUrl(fid, *g_DownloadUrl)

	retStr = fmt.Sprintf(`{"code":0, "data":{"url":"%s"}}`, url)

	logger.Info("return:", retStr)

	_, err := ctx.Write([]byte(retStr))

	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}
