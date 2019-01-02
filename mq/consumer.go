package mq

import "C"
import (
	"context"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

// HandlerFunc is a function that handles message callbacks
type HandlerFunc func(context.Context, *amqp.Delivery)

// Consumer is RabbitMQ consumer
type Consumer struct {
	Ctx         context.Context
	Name        string
	Parallelism int
	Handler     HandlerFunc
	Queue       string
	AutoAck     bool

	Connection  *Connection
	Channels    []*amqp.Channel
	isConnected bool
}

// NewConsumer return a consumer client.
func NewConsumer(ctx context.Context, name, queue string, baseConn *Connection, parallelism int, autoAct bool, handler HandlerFunc) (*Consumer, error) {

	consumer := &Consumer{
		Ctx:         ctx,
		Connection:  baseConn,
		Name:        name,
		Channels:    make([]*amqp.Channel, 0, parallelism),
		Parallelism: parallelism,
		Handler:     handler,
		Queue:       queue,
		AutoAck:     autoAct,
	}

	err := consumer.connect()

	return consumer, err
}

func (cm *Consumer) connect() (err error) {

	var i = 0

	defer func() {
		if err != nil {
			for ; i > 0; i-- {
				cm.Channels[i-1].Close()
			}
		}
	}()

	for ; i < cm.Parallelism; i++ {
		var ch *amqp.Channel
		ch, err = cm.Connection.Channel()
		if err != nil {
			return
		}

		if err = ch.Confirm(false); err != nil {
			return
		}

		if _, err = ch.QueueDeclare(
			cm.Queue,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return
		}

		cm.Channels = append(cm.Channels, ch)
	}

	cm.isConnected = true

	return
}

// Close a consumer client.
func (cm *Consumer) Close() error {
	var err error

	if !cm.isConnected {
		return errors.New("mq consumer have not been init")
	}

	for _, v := range cm.Channels {
		err = v.Close()
		if err != nil {
			log.Error("close consumer(%s) failed, err: ", err)
		}
	}

	return err
}

func (cm *Consumer) consumerProc(ctx context.Context, name string, channel *amqp.Channel) {
	var err error

	if !cm.isConnected {
		log.Error("consumer have not been init")
		return
	}

	var msgs <-chan amqp.Delivery

	if msgs, err = channel.Consume(cm.Queue,
		"",
		cm.AutoAck,
		false,
		false,
		false,
		nil); err != nil {

		log.Error("consumerProc failed to open a Consume, err: ", err)
		return
	}

	for {
		select {
		case msg, ok := <-msgs:
			if ok {
				log.Debugf("consumer id: %s\n", name)
				cm.Handler(ctx, &msg)
			}
		case <-ctx.Done():
			log.Debug(name, "监控退出，停止了...")
			return
		}
	}
}

// Start all goroutines on the consumer client
func (cm *Consumer) Start() error {

	if !cm.isConnected {
		return log.Error("have not connection")
	}

	for i := range cm.Channels {
		go cm.consumerProc(cm.Ctx, fmt.Sprintf("%s gorouine [%d]", cm.Name, i), cm.Channels[i])
	}

	return nil
}
