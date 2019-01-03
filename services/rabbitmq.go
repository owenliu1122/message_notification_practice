package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
)

/*Connections*/

// NewMQConnection return a rabbitmq connection.
func NewMQConnection(url string) (*Connection, error) {
	var err error
	conn := &Connection{
		URL: url,
	}
	err = conn.connect()
	return conn, err
}

// Connection is a rabbitmq connection, it include rabbitmq url address.
type Connection struct {
	URL        string
	Connection *amqp.Connection
}

func (conn *Connection) connect() (err error) {
	conn.Connection, err = amqp.Dial(conn.URL)
	if err != nil {
		return
	}

	return
}

// Channel return a *amqp.Channel.
func (conn *Connection) Channel() (*amqp.Channel, error) {
	return conn.Connection.Channel()
}

// Close current amqp connection.
func (conn *Connection) Close() error {
	return conn.Connection.Close()
}

/*Consumer*/

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

/*Producer*/

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
