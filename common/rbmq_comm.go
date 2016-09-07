package common

//
//import (
//	"fmt"
//
//	"github.com/donnie4w/go-logger/logger"
//	"github.com/streadway/amqp"
//)
//
//type RbmqPublishItem struct {
//	Exchange string
//	Key      string
//	Msg      string
//}
//
//var g_rbmq *amqp.Channel = nil
//var g_itemch chan RbmqPublishItem
//
//func RbmqInit(uri string) int {
//
//	logger.Info("Uri: ", uri)
//
//	conn, err := amqp.Dial(uri)
//	if err != nil {
//		logger.Info("Dial: ", uri, " fail:", err)
//		return -1
//	}
//
//	g_rbmq, err = conn.Channel()
//	if err != nil {
//		logger.Error("Channel:  ", err.Error())
//		return -1
//	}
//
//	g_itemch = make(chan RbmqPublishItem, 1000)
//	go RbmqProduce(g_itemch)
//
//	return 0
//}
//
//func RbmqDeclareExchange(name string, kind string) int {
//
//	err := g_rbmq.ExchangeDeclare(
//		name,  // name
//		kind,  // type
//		false, // durable
//		false, // auto-deleted
//		false, // internal
//		false, // noWait
//		nil,   // arguments
//	)
//
//	if err != nil {
//		logger.Error("ExchangeDeclare fail :", err.Error())
//		return -1
//	}
//
//	return 0
//}
//
//func RbmqProduce(itemch chan RbmqPublishItem) {
//
//	for {
//		select {
//		case item := <-itemch:
//
//			err := g_rbmq.Publish(
//				item.Exchange, // exchange
//				item.Key,      // routing key
//				false,         // mandatory
//				false,         // immediate
//				amqp.Publishing{
//					ContentType: "text/plain",
//					Body:        []byte(item.Msg),
//				})
//
//			if err != nil {
//				logger.Error("sync to data center fail, ", err)
//			} else {
//				logger.Info("sync to data center suc")
//			}
//		}
//	}
//
//	return
//}
//
//func (user *UserInfo) SyncRegToDataCenter(exchange string, key string) {
//
//	msg := "{\"type\":\"Reg\",\"data\":["
//	msg += fmt.Sprintf(`{"app_key":"%s","uid":%d,"cid":"%s","platform":"%s","reg_date":%d}`, user.Appkey, user.Uid, user.Cid, user.Platform, user.RegDate)
//	msg += "]}"
//
//	logger.Info("msg to rbmq: ", msg)
//
//	item := RbmqPublishItem{exchange, key, msg}
//	g_itemch <- item
//
//	return
//}
//
//func (send *MsgSendItem) SyncSendMsgToDataCenter(exchange string, key string) {
//
//	msg := "{\"type\":\"Msg_Send\",\"data\":["
//	msg += fmt.Sprintf(`{"msg_id":%d,"from_uid":%d,"to_uid":%d,"content":"%s","send_time":%d,"type":%d}`, send.Msgid, send.Fromuid, send.Touid, send.Content, send.SendTime, send.Type)
//	msg += "]}"
//
//	logger.Info("msg to rbmq: ", msg)
//
//	item := RbmqPublishItem{exchange, key, msg}
//	g_itemch <- item
//
//	return
//}
//
//func (recv *MsgRecvItem) SyncRecvMsgToDataCenter(exchange string, key string) {
//
//	msg := "{\"type\":\"Msg_Recv\",\"data\":["
//	msg += fmt.Sprintf(`{"msg_id":%d,"uid":%d,"recv_time":%d,"type":%d}`, recv.Msgid, recv.Uid, recv.RecvTime, recv.Type)
//	msg += "]}"
//
//	logger.Info("msg to rbmq: ", msg)
//
//	item := RbmqPublishItem{exchange, key, msg}
//	g_itemch <- item
//
//	return
//}
