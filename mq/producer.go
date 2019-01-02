package mq

import (
	"errors"
	"time"

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

// DeclareExpiration to declare delay queue function.
func (pc *Producer) DeclareExpiration(exchange, routingKey, delayExchange, delayRouting string, expiration time.Duration) error {

	/**
	* 注意,这里是重点!!!!!
	* 声明一个延时队列, 我们的延时消息就是要发送到这里
	 */
	//delayRouting := routingKey
	//delayExchange := exchange
	//delayRouting := routingKey + "_middle_delay"
	//delayExchange := exchange + "_middle_delay"
	delayQueue := routingKey

	if err := pc.Channel.ExchangeDeclare(exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil); err != nil {
		return log.Errorf("Failed to declare a delay_exchange, err:", err)
	}

	if err := pc.Channel.ExchangeDeclare(delayExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil); err != nil {
		return log.Errorf("Failed to declare a delay_exchange, err:", err)
	}

	if _, errDelay := pc.Channel.QueueDeclare(
		delayQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		amqp.Table{
			// 当消息过期时把消息发送到 logs 这个 exchange
			"x-message-ttl":             int32(formatMs(expiration)),
			"x-dead-letter-exchange":    delayExchange,
			"x-dead-letter-routing-key": delayRouting,
		}, // arguments
	); errDelay != nil {
		return log.Errorf("Failed to declare a delay_queue, err:", errDelay)
	}

	if errBind := pc.Channel.QueueBind(
		delayQueue, // queue name, 这里指的是 test_logs
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil); errBind != nil {
		return log.Errorf("Failed to bind a delay_queue, err:", errBind)
	}

	return nil
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

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		log.Errorf(
			"specified duration is %s, but minimal supported value is %s",
			dur, time.Millisecond,
		)
	}
	return int64(dur / time.Millisecond)
}
