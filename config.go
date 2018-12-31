package notice

// Config is all service types configurations.
type Config struct {
	Dashboard    Dashboard    `json:"dashboard" yaml:"dashboard"`
	Server       Server       `json:"server" yaml:"server"`
	Notification Notification `json:"notification" yaml:"notification"`
	Sender       Sender       `json:"sender" yaml:"sender"`
}

// Dashboard is dashboard service processing used configuration.
type Dashboard struct {
	MySQL string `json:"mysql" yaml:"mysql"`
	Redis string `json:"redis" yaml:"redis"`
}

// Server is server service processing used configuration.
type Server struct {
	MySQL    string   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string   `json:"redis" yaml:"redis"`
	RabbitMQ string   `json:"rabbitmq" yaml:"rabbitmq"`
	Producer Producer `json:"producer" yaml:"producer"`
}

// Notification is notification service processing used configuration.
type Notification struct {
	MySQL    string              `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis    string              `json:"redis" yaml:"redis"`
	RabbitMQ string              `json:"rabbitmq" yaml:"rabbitmq"`
	Consumer Consumer            `json:"consumer" yaml:"consumer"`
	Producer map[string]Producer `json:"producer" yaml:"producer"`
}

// Sender is sender service processing used configuration.
type Sender struct {
	MySQL       string              `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis       string              `json:"redis" yaml:"redis"`
	RabbitMQ    string              `json:"rabbitmq" yaml:"rabbitmq"`
	SendService SendService         `json:"sendservice" yaml:"sendservice"`
	Consumer    map[string]Consumer `json:"consumer" yaml:"consumer"`
	Producer    map[string]Producer `json:"producer" yaml:"producer"`
}

// SendService is sender service API configurations.
type SendService struct {
	// mail sender service configurations
	Domain        string `json:"domain" yaml:"domain"`
	PrivateAPIKey string `json:"privateapikey" yaml:"privateapikey"`
	PublicAPIKey  string `json:"publicapikey" yaml:"publicapikey"`

	// phone sender service configurations

	// wechat sender service configurations
}

// Consumer is mq Consumer config information struct.
type Consumer struct {
	Queue string `json:"queue" yaml:"queue"`
}

// Producer is mq producer config information struct.
type Producer struct {
	Exchange   string `json:"exchange" yaml:"exchange"`
	RoutingKey string `json:"routingkey" yaml:"routingkey"`
}
