package common

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/streadway/amqp"
)

type MqMsg struct {
	Msg  []byte
	Rkey string
}

type MqItem struct {
	uri                 string
	ex, rk, queue       string
	readCh              chan []byte
	writeCh             chan MqMsg
	readQuit, writeQuit chan error
	stopRead, stopWrite chan error
	conn                *amqp.Connection
	channel             *amqp.Channel
}

func NewMqItem(uri, ex, rk, queue string, readCh chan []byte, writeCh chan MqMsg) *MqItem {
	readQuit := make(chan error)
	writeQuit := make(chan error)
	stopRead := make(chan error)
	stopWrite := make(chan error)

	return &MqItem{uri, ex, rk, queue, readCh, writeCh, readQuit, writeQuit, stopRead, stopWrite, nil, nil}
}

func (pi *MqItem) clear() {
	time.Sleep(time.Second * 5)

	select {
	case <-pi.readQuit:
		logger.Info("clear readQuit")
	case <-time.After(time.Second):
	}

	select {
	case <-pi.writeQuit:
		logger.Info("clear writeQuit")
	case <-time.After(time.Second):

	}

	select {
	case <-pi.stopRead:
		logger.Info("clear stopRead")
	case <-time.After(time.Second):

	}

	select {
	case <-pi.stopWrite:
		logger.Info("clear stopWrite")
	case <-time.After(time.Second):

	}

	if pi.conn != nil {
		pi.conn.Close()
	}

	if pi.channel != nil {
		pi.channel.Close()
	}

	logger.Info("Clear ok")
}

func (pi *MqItem) Start() int {
	for {
		var err error
		pi.conn, err = amqp.Dial(pi.uri)
		if err != nil {
			logger.Info("apns mq dial: ", err.Error())
			continue
		}

		pi.channel, err = pi.conn.Channel()

		if err != nil {
			logger.Info("apns mq channel: ", err.Error())
			continue
		}

		err = pi.channel.ExchangeDeclare(
			pi.ex,    // name
			"direct", // type
			false,    // durable
			false,    // auto-deleted
			false,    // internal
			false,    // noWait
			nil,      // arguments
		)

		if err != nil {
			logger.Info("apns mq channel: ", err.Error())
			continue
		}

		go func(done chan error) {
			defer func() {
				pi.writeQuit <- nil
			}()

			for {
				select {
				case msg := <-pi.writeCh:
					logger.Debug("gameresp: \n", hex.Dump(msg.Msg))
					if err = pi.channel.Publish(
						pi.ex,    // publish to an exchange
						msg.Rkey, // routing to 0 or more queues
						false,    // mandatory
						false,    // immediate
						amqp.Publishing{
							ContentType: "application/octet-stream",
							Body:        msg.Msg,
						},
					); err != nil {
						logger.Error("Exchange Publish: %s", err)
						return
					} else {
						logger.Debug("Exchange Publish to ", pi.uri, " OK")
					}
				case <-done:
					return
				}
			}
		}(pi.stopWrite)

		if len(pi.queue) > 1 && pi.writeCh != nil {

			queue, err := pi.channel.QueueDeclare(
				pi.queue, // name of the queue
				false,    // durable
				false,    // delete when usused
				false,    // exclusive
				false,    // noWait
				nil,      // arguments
			)

			if err != nil {
				pi.stopWrite <- nil
				logger.Info("apns mq queue: ", err.Error())
				continue
			}

			logstr := fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange ",
				queue.Name, queue.Messages, queue.Consumers)

			logger.Info(logstr)

			if err = pi.channel.QueueBind(
				queue.Name, // name of the queue
				pi.rk,      // bindingKey
				pi.ex,      // sourceExchange
				false,      // noWait
				nil,        // arguments
			); err != nil {
				pi.stopWrite <- nil
				logger.Error("Queue Bind err ", err.Error())
				continue
			}

			deliveries, err := pi.channel.Consume(
				queue.Name, // name
				queue.Name, // consumerTag,
				true,       // noAck
				false,      // exclusive
				false,      // noLocal
				false,      // noWait
				nil,        // arguments
			)

			if err != nil {
				pi.stopWrite <- nil
				logger.Error("Queue Consume: ", err)
				continue
			}

			go pi.handle(deliveries)
		}

		select {
		case <-pi.writeQuit:
			logger.Info(pi.uri, "write proc quit")
		case <-pi.readQuit:
			logger.Info(pi.uri, "read proc quit")
			pi.stopWrite <- nil
		}

		pi.clear()
		time.Sleep(time.Second * 30)
	}
}

func (pi *MqItem) handle(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		pi.readCh <- []byte(d.Body)
	}

	logger.Info("handle: deliveries channel closed")
	pi.readQuit <- nil
}

func (pi *MqItem) WriteMsg(msg MqMsg) {
	pi.writeCh <- msg
}

type MqManager struct {
	uris, ex, rk, queue string
	readCh              chan []byte
	writeCh             chan MqMsg
}

func NewMqManager(uris, ex, rk, queue string, readCh chan []byte) *MqManager {
	writeCh := make(chan MqMsg)
	return &MqManager{uris, ex, rk, queue, readCh, writeCh}
}

func (mm *MqManager) Start() {
	urisArr := strings.Split(mm.uris, ",")

	for _, uri := range urisArr {
		it := NewMqItem(uri, mm.ex, mm.rk, mm.queue, mm.readCh, mm.writeCh)
		go it.Start()
	}

	select {}
}

func (mm *MqManager) Write(msg []byte, rkey string) {
	mm.writeCh <- MqMsg{msg, rkey}
}
