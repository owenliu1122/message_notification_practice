package notice

import "time"

// ConsumerInterface is a redis consumer interface
type ConsumerInterface interface {
	Close() error
	Start() error
}

// ProducerInterface is a redis producer interface
type ProducerInterface interface {
	Publish(exchange, routingKey string, body []byte) error
	DeclareExpiration(exchange, routingKey, delayExchange, delayRouting string, expiration time.Duration) error
	Close() error
}
