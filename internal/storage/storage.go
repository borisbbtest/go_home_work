package storage

type Service_short_url struct {
	ChannelPost chan *string
	ChannelGet  chan *string
	Config      Service_short_urlConfig
}

type Service_short_urlConfig struct {
	Port          int    `yaml:"port"`
	QueueCapacity int    `yaml:"queueCapacity"`
	ServerHost    string `yaml:"ServerHost"`
}
type repositories interface {
}
