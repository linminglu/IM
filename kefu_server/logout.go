package kefu

import (
	"fmt"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/hoisie/web"
)

func (h *Handler) LogOut(ctx *web.Context, val string) {
	logger.Debug("LogOut start")

	retStr := ""

	defer func() {
		logger.Info("return :", retStr)
		ctx.Write([]byte(retStr))
	}()

	reqBuf := make([]byte, 1024)
	i := 0
	strLen := 0
	var err error

	for i < 4 {
		i++
		strLen, err = ctx.Request.Body.Read(reqBuf)

		if strLen > 0 {
			break
		}

		if err != nil {
			logger.Info("readlen :", strLen, "read req:", err)
			time.Sleep(time.Second)
			continue
		}
	}

	logger.Debug("req:", string(reqBuf[0:strLen]))

	retCookie, err := ctx.Request.Cookie("JSESSIONID")

	if err == nil {
		delete(h.Session, retCookie.Value)
	}

	logger.Info("logout success")
	retStr = fmt.Sprintf(`{"code":0, "type":1}`)

	return
}
