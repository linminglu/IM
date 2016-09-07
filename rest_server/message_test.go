package rest_server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"sirendaou.com/duserver/common"
)

func TestMessage(t *testing.T) {
	req := SendReq{
		FromUser: 100057,
		ToUsers:  []uint64{100056},
		Message: Msg{
			Type:   "text",
			Msg:    "test",
			Action: "tst",
		},
		Ext: "ext",
	}

	http.Post("http://192.168.20.51:8889/rest/message/send")
}
