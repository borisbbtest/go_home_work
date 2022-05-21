package config

import (
	"fmt"
	"io/ioutil"

	goflag "flag"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var log = logrus.WithField("context", "service_short_url")

type ServiceShortURLConfig struct {
	Port          int    `yaml:"port"`
	ServerHost    string `yaml:"ServerHost"`
	BaseURL       string `yaml:"BASE_URL"`
	ServerAddress string `yaml:"SERVER_ADDRESS"`
	FileStorePath string `yaml:"FILE_STORAGE_PATH"`
}
type ConfigFromENV struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FileStorePath string `env:"FILE_STORAGE_PATH"`
}
type ServerConfig interface {
	GetConfig() (config *ServiceShortURLConfig, err error)
}

func GetConfig() (config *ServiceShortURLConfig, err error) {

	var ServerAddress, BaseURL, FilePath, configFileName string
	flag.StringVarP(&configFileName, "config", "c", "./config.yml", "path to the configuration file")
	flag.StringVarP(&ServerAddress, "server", "a", "", "Server Adders")
	flag.StringVarP(&BaseURL, "base_url", "b", "", "Base URL")
	flag.StringVarP(&FilePath, "file_path", "f", "", "Config file path")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	log.Infof("Loading configuration at '%s'", configFileName)
	configFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Errorf("can't open the config file: %s", err)

	}
	// Default values
	config = &ServiceShortURLConfig{
		Port:          8080,
		ServerHost:    "localhost",
		BaseURL:       "http://localhost:8080",
		ServerAddress: "localhost:8080",
		FileStorePath: "",
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Errorf("can't read the config file: %s", err)
	}

	var cfgenv ConfigFromENV
	err = env.Parse(&cfgenv)
	if err != nil {
		log.Errorf("can't start the listening thread: %s", err)
	} else {
		if cfgenv.ServerAddress != "" {
			config.ServerAddress = cfgenv.ServerAddress
		}
		if cfgenv.BaseURL != "" {
			config.BaseURL = cfgenv.BaseURL
		}
		if cfgenv.FileStorePath != "" {
			config.FileStorePath = cfgenv.FileStorePath
		}
	}

	if ServerAddress != "" {
		config.ServerAddress = ServerAddress
	}
	if BaseURL != "" {
		config.BaseURL = BaseURL
	}
	if FilePath != "" {
		config.FileStorePath = FilePath
	}

	log.Info("Configuration loaded")
	return
}
