package mq

import "github.com/streadway/amqp"

// NewConnection return a rabbitmq connection.
func NewConnection(url string) (*Connection, error) {
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
