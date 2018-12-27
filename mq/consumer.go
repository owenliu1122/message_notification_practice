package mq

import (
	"context"
	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

/*
Consumer
*/
type HandlerFunc func(context.Context, *amqp.Delivery)

type BaseConsumer struct {
	Queue   string
	AutoAck bool
	NoLocal bool
	//NoAck			bool
	Exclusive bool
	NoWait    bool
	Arguments amqp.Table
}

type ConsumerContext struct {
	ChannelKey string
	Channel    *amqp.Channel
	Handler    HandlerFunc
	BaseConsumer
}

func (mq *BaseMq) RegisterConsumer(name string, handler HandlerFunc, baseCm BaseConsumer) error {

	if mq.cm == nil {
		mq.cm = make(map[string]*ConsumerContext)
	}

	mq.cm[name] = &ConsumerContext{
		ChannelKey:   name,
		Handler:      handler,
		BaseConsumer: baseCm,
	}

	return nil
}

func (mq *BaseMq) closeConsumer() error {
	var err error

	if mq.cm == nil || len(mq.cm) == 0 {
		return err
	}

	for k, v := range mq.cm {
		err = v.Channel.Close()
		if err != nil {
			log.Error("close consumer(%s) failed, err: ", k, err)
		}
	}

	return err
}

func (mq *BaseMq) consumerProc(name string, ctx context.Context, cmCtx *ConsumerContext) {
	var err error

	if mq.conn == nil {
		if err := mq.refreshMqConnection(); err != nil {
			log.Error("consumerProc refresh connection failed, err: ", err)
			return
		}
	}

	if mq.cm[name].Channel, err = mq.conn.Connection.Channel(); err != nil {
		log.Error("consumerProc failed to open a channel, err: ", err)
		return
	}

	var msgs <-chan amqp.Delivery

	if msgs, err = mq.cm[name].Channel.Consume(mq.cm[name].Queue,
		"",
		mq.cm[name].AutoAck,
		mq.cm[name].Exclusive,
		mq.cm[name].NoLocal,
		mq.cm[name].NoWait,
		mq.cm[name].Arguments); err != nil {

		log.Error("consumerProc failed to open a Consume, err: ", err)
		return
	}

	for {
		select {
		case msg, ok := <-msgs:
			if ok {
				log.Debugf("consumer id: %s\n", name)
				mq.cm[name].Handler(ctx, &msg)
			}
		case <-ctx.Done():
			log.Debug(name, "监控退出，停止了...")
			return
		}
	}
}

func (mq *BaseMq) StartConsumer(ctx context.Context) error {

	var err error

	if mq.cm == nil || len(mq.cm) == 0 {
		return err
	}

	for k, v := range mq.cm {
		go mq.consumerProc(k, ctx, v)
	}

	return err
}
