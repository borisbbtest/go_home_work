package storage

type service_short_url struct {
	channelPost chan *string
	channelGet  chan *string
	config      service_short_urlConfig
}

type service_short_urlConfig struct {
	Port          int    `yaml:"port"`
	QueueCapacity int    `yaml:"queueCapacity"`
	ServerHost    string `yaml:"ServerHost"`
}
