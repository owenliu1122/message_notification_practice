package mq

import (
	"errors"

	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

// Producer is a producer client.
type Producer struct {
	Name string
	//Exchange   string
	//RoutingKey string
	//Mandatory  bool
	//Immediate  bool

	Connection  *Connection
	Channel     *amqp.Channel
	isConnected bool
}

// NewProducer return a producer client.
func NewProducer(name string, baseConn *Connection) (*Producer, error) {
	var err error
	producer := &Producer{
		Name:       name,
		Connection: baseConn,
	}
	err = producer.connect()
	return producer, err
}

func (pc *Producer) connect() (err error) {

	var ch *amqp.Channel
	ch, err = pc.Connection.Channel()
	if err != nil {
		return
	}

	if err = ch.Confirm(false); err != nil {
		return
	}

	pc.Channel = ch
	pc.isConnected = true

	return
}

// Publish a message.
func (pc *Producer) Publish(exchange, routingKey string, body []byte) error {

	if !pc.isConnected {
		return errors.New("MQ producer has not been initialized")
	}

	if err := pc.Channel.ExchangeDeclare(exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil); err != nil {
		return err
	}

	return pc.Channel.Publish(exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

// Close a producer client.
func (pc *Producer) Close() error {
	var err error

	if !pc.isConnected {
		return errors.New("mq producer have not been init")
	}

	err = pc.Channel.Close()
	if err != nil {
		return log.Error("close producer(%s) failed, err: ", err)
	}

	return err
}
