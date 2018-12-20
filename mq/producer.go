package mq

import (
	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

func producerProc(id int, info MQInfo) {
	var conn *amqp.Connection
	var channel *amqp.Channel
	//var count = 0
	var err error

	conn, err = amqp.Dial(info.Cfg.URL)
	failOnErr(err, "failed to connect tp rabbitmq")
	defer conn.Close()

	channel, err = conn.Channel()
	failOnErr(err, "failed to open a channel")
	defer channel.Close()

	for {
		select {
		case msg, ok := <-info.MsgChan:
			if ok {
				log.Debugf("producer id: %d\n", id)

				channel.Publish(info.Cfg.Exchange, info.Cfg.Queue, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        msg.([]byte),
				})
			}
		}
	}
}

// 启动消费协程
func ProducerStart(routineNum int, info MQInfo) error {

	for i := 0; i < routineNum; i++ {
		go producerProc(i, info)
	}

	return nil
}
