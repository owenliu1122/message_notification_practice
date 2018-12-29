package mq

import (
	"github.com/streadway/amqp"
	"sync"
)

/*
BaseMq
*/
type mqConnection struct {
	Lock       sync.RWMutex
	Connection *amqp.Connection
	MqURL      string
}

type BaseMq struct {
	conn *mqConnection
	cm   map[string]*ConsumerContext
	pc   *ProducerContext
}

func NewMq(url string) *BaseMq {
	return &BaseMq{
		conn: &mqConnection{
			MqURL: url,
		},
	}
}

func (mq *BaseMq) InitConnection() error {
	return mq.refreshMqConnection()
}

func (mq *BaseMq) refreshMqConnection() error {
	var err error
	mq.conn.Lock.Lock()
	defer mq.conn.Lock.Unlock()

	if mq.conn.Connection == nil {
		for i := 0; i < 5; i++ {
			mq.conn.Connection, err = amqp.Dial(mq.conn.MqURL)
			if err == nil {
				break
			}
		}
	}

	return err
}

func (mq *BaseMq) Close() error {
	mq.closeProducer()

	mq.closeConsumer()

	// 关闭 MQ 连接
	if mq.conn != nil && mq.conn.Connection != nil {
		mq.conn.Connection.Close()
	}

	return nil
}
