package notice

import (
	"github.com/fpay/foundation-go/database"
	"github.com/fpay/foundation-go/job"
	"github.com/fpay/foundation-go/log"
)

// Config is all service types configurations.
type Config struct {
	Dashboard    DashboardConfig    `mapstructure:"dashboard" json:"dashboard" yaml:"dashboard"`
	Server       ServerConfig       `mapstructure:"server" json:"server" yaml:"server"`
	Notification NotificationConfig `mapstructure:"notification" json:"notification" yaml:"notification"`
	Sender       SenderConfig       `mapstructure:"sender" json:"sender" yaml:"sender"`
}

// DashboardConfig is dashboard service processing used configuration.
type DashboardConfig struct {
	Logger log.Options              `mapstructure:"logger" json:"logger" yaml:"logger"`
	MySQL  database.DatabaseOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  string                   `mapstructure:"redis" json:"redis" yaml:"redis"`
}

// ServerConfig is server service processing used configuration.
type ServerConfig struct {
	Logger   log.Options              `mapstructure:"logger" json:"logger" yaml:"logger"`
	MySQL    database.DatabaseOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string                   `mapstructure:"redis" json:"redis" yaml:"redis"`
	RabbitMQ job.JobManagerOptions    `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	Producer JobConfig                `mapstructure:"producer" json:"producer" yaml:"producer"`
}

// NotificationConfig is notification service processing used configuration.
type NotificationConfig struct {
	Logger   log.Options              `mapstructure:"logger" json:"logger" yaml:"logger"`
	MySQL    database.DatabaseOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string                   `mapstructure:"redis" json:"redis" yaml:"redis"`
	RabbitMQ job.JobManagerOptions    `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	Consumer ConsumerConfig           `mapstructure:"consumer" json:"consumer" yaml:"consumer"`
	Producer map[string]JobConfig     `mapstructure:"producer" json:"producer" yaml:"producer"`
}

// SenderConfig is sender service processing used configuration.
type SenderConfig struct {
	Logger      log.Options              `mapstructure:"logger" json:"logger" yaml:"logger"`
	MySQL       database.DatabaseOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis       string                   `mapstructure:"redis" json:"redis" yaml:"redis"`
	RabbitMQ    job.JobManagerOptions    `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	SendService SendServiceConfig        `mapstructure:"sendservice" json:"sendservice" yaml:"sendservice"`
	Consumer    map[string]JobConfig     `mapstructure:"consumer" json:"consumer" yaml:"consumer"`
	RetryDelay  int                      `mapstructure:"retrydelay" json:"retrydelay" yaml:"retrydelay"`
}

// SendServiceConfig is sender service API configurations.
type SendServiceConfig struct {
	// mail sender service configurations
	Domain        string `mapstructure:"domain" json:"domain" yaml:"domain"`
	PrivateAPIKey string `mapstructure:"privateapikey" json:"privateapikey" yaml:"privateapikey"`
	PublicAPIKey  string `mapstructure:"publicapikey" json:"publicapikey" yaml:"publicapikey"`

	// TODO: phone sender service configurations

	// TODO: wechat sender service configurations
}

// ConsumerConfig is mq Consumer config information struct.
type ConsumerConfig struct {
	Queue string `mapstructure:"queue" json:"queue" yaml:"queue"`
}

// ProducerConfig is mq producer config information struct.
type JobConfig struct {
	Delay int    `mapstructure:"delay" json:"delay" yaml:"delay"`
	Queue string `mapstructure:"queue" json:"queue" yaml:"queue"`
}
