package notice

// Config is all service types configurations.
type Config struct {
	Dashboard    DashboardConfig    `json:"dashboard" yaml:"dashboard"`
	Server       ServerConfig       `json:"server" yaml:"server"`
	Notification NotificationConfig `json:"notification" yaml:"notification"`
	Sender       SenderConfig       `json:"sender" yaml:"sender"`
}

// Dashboard is dashboard service processing used configuration.
type DashboardConfig struct {
	MySQL string `json:"mysql" yaml:"mysql"`
	Redis string `json:"redis" yaml:"redis"`
}

// Server is server service processing used configuration.
type ServerConfig struct {
	MySQL    string         `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string         `json:"redis" yaml:"redis"`
	RabbitMQ string         `json:"rabbitmq" yaml:"rabbitmq"`
	Producer ProducerConfig `json:"producer" yaml:"producer"`
}

// Notification is notification service processing used configuration.
type NotificationConfig struct {
	MySQL    string                    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string                    `json:"redis" yaml:"redis"`
	RabbitMQ string                    `json:"rabbitmq" yaml:"rabbitmq"`
	Consumer ConsumerConfig            `json:"consumer" yaml:"consumer"`
	Producer map[string]ProducerConfig `json:"producer" yaml:"producer"`
}

// Sender is sender service processing used configuration.
type SenderConfig struct {
	MySQL         string                    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis         string                    `json:"redis" yaml:"redis"`
	RabbitMQ      string                    `json:"rabbitmq" yaml:"rabbitmq"`
	RetryDelay    int64                     `json:"retrydelay" yaml:"retrydelay"`
	SendService   SendServiceConfig         `json:"sendservice" yaml:"sendservice"`
	Consumer      map[string]ConsumerConfig `json:"consumer" yaml:"consumer"`
	RetryProducer map[string]ProducerConfig `json:"retryproducer" yaml:"retryproducer"`
	DelayProducer map[string]ProducerConfig `json:"delayproducer" yaml:"delayproducer"`
}

// SendService is sender service API configurations.
type SendServiceConfig struct {
	// mail sender service configurations
	Domain        string `json:"domain" yaml:"domain"`
	PrivateAPIKey string `json:"privateapikey" yaml:"privateapikey"`
	PublicAPIKey  string `json:"publicapikey" yaml:"publicapikey"`

	// TODO: phone sender service configurations

	// TODO: wechat sender service configurations
}

// Consumer is mq Consumer config information struct.
type ConsumerConfig struct {
	Queue string `json:"queue" yaml:"queue"`
}

// Producer is mq producer config information struct.
type ProducerConfig struct {
	Exchange   string `json:"exchange" yaml:"exchange"`
	RoutingKey string `json:"routingkey" yaml:"routingkey"`
}
