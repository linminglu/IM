package rest_server

import (
	"net/http"
	"strings"
)

func (h *Handler) GetSetupID(res http.ResponseWriter, req *http.Request) {
	newSetupID, _ := h.redisMgr.RedisSetupIDCacheGet()
	s := []string {`{"setupid":`, `}`}
	str := strings.Join(s, newSetupID)

	res.Write([]byte(str))
}