package apns_server

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"sirendaou.com/duserver/apns_server/apns"
	"sirendaou.com/duserver/common"
	"time"
)

type ApnsProc struct {
	//	AppKey         string
	Stat           int
	MsgCh          chan common.ApnsMsg
	KeyFilePath    string
	QuitCh         chan string
	UserRedis      *common.RedisManager
	FeedbackClient *apns.Client
	ProcKey        string
}

func NewApnsProc(path string, stat int, procStopCh chan string, userRedis *common.RedisManager, procKey string) *ApnsProc {
	return &ApnsProc {
		Stat:stat,
		MsgCh:make(chan common.ApnsMsg, 1000),
		KeyFilePath:path,
		QuitCh:procStopCh,
		UserRedis:userRedis,
		ProcKey:procKey,
	}
}

func (p *ApnsProc) isOnline(uid uint64) bool {
	_, err := p.UserRedis.RedisStatCacheGet(uid)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func FeedBackProc(UserRedis *common.RedisManager) {
	for {
		select {
		case resp := <-apns.FeedbackChannel:
			logger.Info("recv token:", resp.DeviceToken, " is invalide")
			uid := common.DBGetUidByToken(resp.DeviceToken)
			if uid > 0 {
				key := fmt.Sprintf("token_%d", uid)
				UserRedis.RedisDel(key)
				logger.Info("feedback clear invalied token:", resp.DeviceToken, "uid:", uid)
			}
		}
	}
}

func (p *ApnsProc) FeedBackStart() {
	keys := ""
	pem := ""
	url := ""
	if p.Stat == 2 {
		//		keys = fmt.Sprintf("%s/%s.key", *g_PemFilePath, p.AppKey)
		//		pem = fmt.Sprintf("%s/%s.pem", *g_PemFilePath, p.AppKey)
		keys = fmt.Sprintf("%s/%s.key", *g_PemFilePath, "du")
		pem = fmt.Sprintf("%s/%s.pem", *g_PemFilePath, "du")
		//url = "gateway.push.apple.com:2195"
		url = "feedback.push.apple.com:2196"
	} else {
		//		keys = fmt.Sprintf("%s/%s_test.key", *g_PemFilePath, p.AppKey)
		//		pem = fmt.Sprintf("%s/%s_test.pem", *g_PemFilePath, p.AppKey)
		keys = fmt.Sprintf("%s/%s_test.key", *g_PemFilePath, "du")
		pem = fmt.Sprintf("%s/%s_test.pem", *g_PemFilePath, "du")
		//url = "gateway.sandbox.push.apple.com:2195"
		url = "feedback.sandbox.push.apple.com:2196"
	}

	for {
		p.FeedbackClient = apns.NewClient(url, pem, keys)
		err, ret := p.FeedbackClient.ListenForFeedback()

		logger.Info("ListenForFeedback err ", err, ret)

		if ret == 4 {
			continue
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}

}

func (p *ApnsProc) FeedBackStop() {
	return
	//p.FeedbackClient.FeedBackquitch <- 0
}

func (p *ApnsProc) StartProc() {
	keys := ""
	pem := ""
	url := ""

	if p.Stat == 2 {
		keys = fmt.Sprintf("%s/%s.key", *g_PemFilePath, "du")
		pem = fmt.Sprintf("%s/%s.pem", *g_PemFilePath, "du")
		url = "gateway.push.apple.com:2195"
	} else {
		keys = fmt.Sprintf("%s/%s_test.key", *g_PemFilePath, "du")
		pem = fmt.Sprintf("%s/%s_test.pem", *g_PemFilePath, "du")
		url = "gateway.sandbox.push.apple.com:2195"
	}

	p.QuitCh <- p.ProcKey

	for {
		client := apns.NewClient(url, pem, keys)

		err := client.Connect()
		if err != nil {
			logger.Info(" connect to ", url, " fail:", err)
			return
		} else {
			logger.Info(" connect to ", url, " success")
		}

		defer func() {
			client.Close()
		}()

		logger.Info("start apns server:", url)

		payload := apns.NewPayload()
		payload.Alert = ""
		payload.Badge = 1
		payload.Sound = "bingbong.aiff"

		isContinue := true

		for isContinue {
			select {
			case msg := <-p.MsgCh:
				logger.Info("recv msg :", msg)
				payload.Alert = msg.Msg
				var pipeKeys []string
				token := ""
				for _, uid := range msg.UidList {
					if p.isOnline(uid) {
						logger.Info(uid, " is online")
					} else {
//						_, pt := common.GetAppIdPtFromUid(uid)
//						if pt == common.PT_IOS {
							token = fmt.Sprintf("token_%d", uid)
							pipeKeys = append(pipeKeys, token)
//						} else {
//							logger.Debug("Uid:", uid, "is not ios")
//						}
					}
				}

				tokenList := p.UserRedis.PipelineGetString(pipeKeys)

				logger.Debug("tokenlist ===", tokenList)

				for _, token := range tokenList {
					logger.Debug("send msg to token:", token)

					if len(token) < 60 {
						break
					}

					pn := apns.NewPushNotification()
					pn.DeviceToken = token
					pn.AddPayload(payload)

					// send to apns
					resp := client.Send(pn)
					pn.PayloadString()

					logger.Info("Success:", resp.Success)

					if resp.Error != nil {
						logger.Error("------Error:", resp.Error)
						isContinue = false
						break
					}
				}
			}
		}

		logger.Info("reconnect apns server")
	}
}

func (p *ApnsProc) sendMsg(apnsReq common.ApnsMsg) int {
	logger.Debug("add msg :", apnsReq)
	p.MsgCh <- apnsReq
	return 0
}
