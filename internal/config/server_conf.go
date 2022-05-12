package config

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log = logrus.WithField("context", "service_short_url")

type ServiceShortURLConfig struct {
	Port           int    `yaml:"port"`
	ServerHost     string `yaml:"ServerHost"`
	BASE_URL       string `yaml:"BASE_URL"`
	SERVER_ADDRESS string `yaml:"SERVER_ADDRESS"`
}

type ConfigFromENV struct {
	SERVER_ADDRESS string `env:"SERVER_ADDRESS,required"`
	BASE_URL       string `env:"BASE_URL,required"`
}
type ServerConfig interface {
	getConfig(filename string) (cfg *ServiceShortURLConfig, err error)
}

func GetConfig(filename string) (cfg *ServiceShortURLConfig, err error) {
	log.Infof("Loading configuration at '%s'", filename)
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open the config file: %s", err)
	}

	// Default values
	config := ServiceShortURLConfig{
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
