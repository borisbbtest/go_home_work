package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.WithField("context", "service_short_url")

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

func New(cfg *service_short_urlConfig) *service_short_url {

	return &service_short_url{
		channelPost: make(chan *string, cfg.QueueCapacity),
		channelGet:  make(chan *string, cfg.QueueCapacity),
		config:      *cfg,
	}
}

func ConfigFromFile(filename string) (cfg *service_short_urlConfig, err error) {
	log.Infof("Loading configuration at '%s'", filename)
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open the config file: %s", err)
	}

	// Default values
	config := service_short_urlConfig{
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

func (hook *service_short_url) Start() error {

	// Launch the process thread
	go hook.processMain()

	// Launch the listening thread
	log.Println("Initializing HTTP server")
	http.HandleFunc("/", hook.mainHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(hook.config.Port), nil)
	if err != nil {
		return fmt.Errorf("can't start the listening thread: %s", err)
	}

	log.Info("Exiting")
	close(hook.channelGet)
	close(hook.channelPost)

	return nil
}

func (hook *service_short_url) mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hook.PostHandler(w, r)
	case http.MethodGet:
		hook.GetHandler(w, r)
	default:
		http.Error(w, "unsupported HTTP method only post send", 400)
	}
}

func (hook *service_short_url) GetHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(bytes))

	defer r.Body.Close()
	var m string
	if err := json.Unmarshal(bytes, &m); err != nil {
		log.Errorf("body error: %v", string(bytes))
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "request body is not valid json", 400)
		return
	}

	hook.channelGet <- &m
	fmt.Printf(m)

}

func (hook *service_short_url) PostHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(bytes))

	defer r.Body.Close()
	var m string
	if err := json.Unmarshal(bytes, &m); err != nil {
		log.Errorf("body error: %v", string(bytes))
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "request body is not valid json", 400)
		return
	}

	hook.channelPost <- &m
	fmt.Printf(m)

}
func (hook *service_short_url) processMain() {
	log.Info("Get URL to short")
	for {
		select {
		case a := <-hook.channelPost:
			if a == nil {
				log.Info("Queue Closed")
				return
			}
		case a := <-hook.channelPost:
			if a == nil {
				log.Info("Queue Closed")
				return
			}
		default:
		}
	}
}
