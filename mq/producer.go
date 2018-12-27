package mq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

/*
Producer
*/

type BaseProducer struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
}

type ProducerContext struct {
	//ChannelKey		string
	BaseProducer
	Channel *amqp.Channel
}

func (mq *BaseMq) InitProducer(exchange, routingKey string) error {
	mq.pc = &ProducerContext{
		BaseProducer: BaseProducer{
			Exchange:   exchange,
			RoutingKey: routingKey,
		},
	}

	return nil
}

func (mq *BaseMq) refreshProducerChannel() error {
	var err error

	err = mq.refreshMqConnection()
	if err != nil {
		return errors.New(fmt.Sprintf("refreshMqConnection failed, err: %s", err.Error()))
	}

	mq.pc.Channel, err = mq.conn.Connection.Channel()

	return err
}

func (mq *BaseMq) Send(exchange, routingKey string, body []byte) error {

	if mq.pc == nil {
		return errors.New("MQ producer has not been initialized")
	}

	if err := mq.refreshProducerChannel(); err != nil {
		return errors.New(fmt.Sprintf("refreshProducerChannel failed, err: %s", err.Error()))
	}

	if "" != exchange {
		mq.pc.Exchange = exchange
	}

	if "" != routingKey {
		mq.pc.RoutingKey = routingKey
	}

	return mq.pc.Channel.Publish(mq.pc.Exchange,
		mq.pc.RoutingKey,
		mq.pc.Mandatory,
		mq.pc.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func (mq *BaseMq) closeProducer() error {
	var err error
	// 关闭生产者 Channel
	if mq.pc != nil && mq.pc.Channel != nil {
		err = mq.pc.Channel.Close()
		if err != nil {
			log.Error("close producer failed, err: ", err)
		}
	}

	return err
}
