package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type ReportReq struct {
	Report_msg string `json:"report_msg,omitempty"`
}

func (h *DBHandler) Report(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req ReportReq
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	if req.Report_msg == "" {
		logger.Error("req.Report_msg is nil")
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	common.DBAddReport(head.Uid, req.Report_msg)

	SendTime := uint32(time.Now().Unix())
	msgId := common.GetChatMsgId()

	resp := fmt.Sprintf(`{"msgid":%d, "sendtime":%d}`, msgId, SendTime)

	return []byte(resp), 0
}
