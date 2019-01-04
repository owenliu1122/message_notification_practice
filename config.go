package notice

import (
	"github.com/fpay/foundation-go/database"
	"github.com/fpay/foundation-go/log"
)

// Config is all service types configurations.
type Config struct {
	Dashboard    DashboardConfig    `json:"dashboard" yaml:"dashboard"`
	Server       ServerConfig       `json:"server" yaml:"server"`
	Notification NotificationConfig `json:"notification" yaml:"notification"`
	Sender       SenderConfig       `json:"sender" yaml:"sender"`
}

// DashboardConfig is dashboard service processing used configuration.
type DashboardConfig struct {
	Logger log.Options              `json:"logger" yaml:"logger"`
	MySQL  database.DatabaseOptions `json:"mysql" yaml:"mysql"`
	Redis  string                   `json:"redis" yaml:"redis"`
}

// ServerConfig is server service processing used configuration.
type ServerConfig struct {
	Logger   log.Options              `json:"logger" yaml:"logger"`
	MySQL    database.DatabaseOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string                   `json:"redis" yaml:"redis"`
	RabbitMQ string                   `json:"rabbitmq" yaml:"rabbitmq"`
	Producer ProducerConfig           `json:"producer" yaml:"producer"`
}

// NotificationConfig is notification service processing used configuration.
type NotificationConfig struct {
	Logger   log.Options               `json:"logger" yaml:"logger"`
	MySQL    database.DatabaseOptions  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string                    `json:"redis" yaml:"redis"`
	RabbitMQ string                    `json:"rabbitmq" yaml:"rabbitmq"`
	Consumer ConsumerConfig            `json:"consumer" yaml:"consumer"`
	Producer map[string]ProducerConfig `json:"producer" yaml:"producer"`
}

// SenderConfig is sender service processing used configuration.
type SenderConfig struct {
	Logger        log.Options               `json:"logger" yaml:"logger"`
	MySQL         database.DatabaseOptions  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis         string                    `json:"redis" yaml:"redis"`
	RabbitMQ      string                    `json:"rabbitmq" yaml:"rabbitmq"`
	RetryDelay    int64                     `json:"retrydelay" yaml:"retrydelay"`
	SendService   SendServiceConfig         `json:"sendservice" yaml:"sendservice"`
	Consumer      map[string]ConsumerConfig `json:"consumer" yaml:"consumer"`
	RetryProducer map[string]ProducerConfig `json:"retryproducer" yaml:"retryproducer"`
	DelayProducer map[string]ProducerConfig `json:"delayproducer" yaml:"delayproducer"`
}

// SendServiceConfig is sender service API configurations.
type SendServiceConfig struct {
	// mail sender service configurations
	Domain        string `json:"domain" yaml:"domain"`
	PrivateAPIKey string `json:"privateapikey" yaml:"privateapikey"`
	PublicAPIKey  string `json:"publicapikey" yaml:"publicapikey"`

	// TODO: phone sender service configurations

	// TODO: wechat sender service configurations
}

// ConsumerConfig is mq Consumer config information struct.
type ConsumerConfig struct {
	Queue string `json:"queue" yaml:"queue"`
}

// ProducerConfig is mq producer config information struct.
type ProducerConfig struct {
	Exchange   string `json:"exchange" yaml:"exchange"`
	RoutingKey string `json:"routingkey" yaml:"routingkey"`
}
