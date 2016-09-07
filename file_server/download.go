package file_server

import (
	"fmt"
	"strconv"

	"github.com/hoisie/web"
 "github.com/donnie4w/go-logger/logger"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/fop"
)

func GetUrl(fid string, url string) string {
	baseUrl := rs.MakeBaseUrl(url, fid)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}

func GetImageUrl(w, h int, fid, url string) string {
	baseUrl := rs.MakeBaseUrl(url, fid)
	var view = fop.ImageView{2, w, h, 80, "jpg"}
	imageUrl := view.MakeRequest(baseUrl)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(imageUrl, nil)
}

func downloadHandlerPost(ctx *web.Context, val string) {
	ctx.Request.ParseMultipartForm(1024)

	fromUid := ctx.Request.FormValue("fromuid")
	fid := ctx.Request.FormValue("fid")
	ver := ctx.Request.FormValue("vercode")

	logger.Info("fid:", fid, "fromuid:", fromUid, "vercode", ver)

	retStr := ""

	for {
		if len(fid) < 32 || len(fromUid) < 1 {
			ctx.NotFound("not found")
			retStr = `{"errcode":1002}`
			break
		}

		if len(ver) != 32 {
			ctx.Unauthorized()
			retStr = `{"errcode":1003}`
			break
		}

		url := GetUrl(fid, *g_DownloadUrl)

		retStr = fmt.Sprintf(`{"errcode":0,"url":"%s"}`, url)
		break
	}

	_, err := ctx.Write([]byte(retStr))
	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}

func downloadHandler(ctx *web.Context, val string) {
	fid := ""
	fromUid := ""
	ver := ""

	for k, v := range ctx.Params {
		if k == "fid" {
			fid = v
		} else if k == "fromuid" {
			fromUid = v
		} else if k == "vercode" {
			ver = v
		}
	}

	logger.Info("fid:", fid, "fromuid:", fromUid, "vercode", ver)

	retStr := ""
	for {
		if len(fid) < 32 || len(fromUid) < 5 {
			ctx.NotFound("not found")
			retStr = `{"errcode":1002}`
			break
		}

		if len(ver) != 32 {
			ctx.Unauthorized()
			retStr = `{"errcode":1003}`
			break
		}

		url := GetUrl(fid, *g_DownloadUrl)

		retStr = fmt.Sprintf(`{"errcode":0,"url":"%s"}`, url)
		break
	}

	_, err := ctx.Write([]byte(retStr))
	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}

func imageHandler(ctx *web.Context, val string) {
	w := 100
	h := 100

	ctx.Request.ParseMultipartForm(1024)

	fromUid := ctx.Request.FormValue("fromuid")
	fid := ctx.Request.FormValue("fid")
	ver := ctx.Request.FormValue("vercode")
	sw := ctx.Request.FormValue("width")
	sh := ctx.Request.FormValue("high")

	w, _ = strconv.Atoi(sw)
	h, _ = strconv.Atoi(sh)

	logger.Info("fid:", fid, "fromuid:", fromUid, "vercode", ver, "w:", w, "h:", h)

	retStr := ""
	for {
		if len(fid) < 32 || len(fromUid) < 5 {
			ctx.NotFound("not found")
			retStr = `{"errcode":1002}`
			break
		}

		if len(ver) != 32 {
			ctx.Unauthorized()
			retStr = `{"errcode":1003}`
			break
		}

		url := GetImageUrl(w, h, fid, *g_DownloadUrl)

		retStr = fmt.Sprintf(`{"errcode":0,"url":"%s"}`, url)
		break
	}

	_, err := ctx.Write([]byte(retStr))
	if err != nil {
		logger.Error("ctx.Write err:", err)
	}

	return
}
