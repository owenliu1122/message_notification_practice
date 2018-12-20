package mq

import (
	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

func consumerProc(id int, info MQInfo) {
	var conn *amqp.Connection
	var channel *amqp.Channel
	//var count = 0
	var err error

	conn, err = amqp.Dial(info.Cfg.URL)
	failOnErr(err, "failed to connect tp rabbitmq")
	defer func() { log.Debug("consumer killed") }()
	defer conn.Close()

	channel, err = conn.Channel()
	failOnErr(err, "failed to open a channel")
	defer channel.Close()

	msgs, err := channel.Consume(info.Cfg.Queue, "", true, false, false, false, nil)
	failOnErr(err, "")

	for {
		select {
		case msg, ok := <-msgs:
			if ok {
				log.Debugf("consumer id: %d\n", id)
				info.MsgChan <- msg
			}

		}
	}
}

// 启动消费协程
func ConsumerStart(routineNum int, info MQInfo) error {

	for i := 0; i < routineNum; i++ {
		go consumerProc(i, info)
	}

	return nil
}
