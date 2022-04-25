package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	. "github.com/borisbbtest/go_home_work/internal/handlers"
	. "github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.WithField("context", "service_short_url")

func New(cfg *Service_short_urlConfig) *Service_short_url {

	return &Service_short_url{
		ChannelPost: make(chan *string, cfg.QueueCapacity),
		ChannelGet:  make(chan *string, cfg.QueueCapacity),
		Config:      *cfg,
	}
}

func ConfigFromFile(filename string) (cfg *Service_short_urlConfig, err error) {
	log.Infof("Loading configuration at '%s'", filename)
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open the config file: %s", err)
	}

	// Default values
	config := Service_short_urlConfig{
		Port:          8080,
		QueueCapacity: 500,
		ServerHost:    "localhost",
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("can't read the config file: %s", err)
	}

	log.Info("Configuration loaded")
	return &config, nil
}

func (hook *Service_short_url) Start() error {

	// Launch the process thread
	go hook.processMain()

	// Launch the listening thread
	log.Println("Initializing HTTP server")
	http.HandleFunc("/", hook.MainHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(hook.config.Port), nil)
	if err != nil {
		return fmt.Errorf("can't start the listening thread: %s", err)
	}

	log.Info("Exiting")
	close(hook.ChannelGet)
	close(hook.ChannelPost)

	return nil
}
func (hook *Service_short_url) ProcessMain() {
	log.Info("Get URL to short")
	for {
		select {
		case a := <-hook.ChannelPost:
			if a == nil {
				log.Info("Queue Closed")
				return
			}
		case a := <-hook.ChannelPost:
			if a == nil {
				log.Info("Queue Closed")
				return
			}
		default:
		}
	}
}
