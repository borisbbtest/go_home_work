package config

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.WithField("context", "service_short_url")

type Service_short_urlConfig struct {
	Port       int    `yaml:"port"`
	ServerHost string `yaml:"ServerHost"`
}
type ServerConfig interface {
	getConfig(filename string) (cfg *Service_short_urlConfig, err error)
}

func GetConfig(filename string) (cfg *Service_short_urlConfig, err error) {
	log.Infof("Loading configuration at '%s'", filename)
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open the config file: %s", err)
	}

	// Default values
	config := Service_short_urlConfig{
		Port:       8080,
		ServerHost: "localhost",
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("can't read the config file: %s", err)
	}

	log.Info("Configuration loaded")
	return &config, nil
}
