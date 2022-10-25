package config

import (
	"encoding/json"
	goflag "flag"
	"io/ioutil"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var log = logrus.WithField("context", "service_short_url")

type ServiceShortURLConfig struct {
	Port          int    `yaml:"port"`
	ServerHost    string `yaml:"ServerHost"`
	BaseURL       string `json:"base_url" yaml:"BASE_URL"`
	ServerAddress string `json:"server_address" yaml:"SERVER_ADDRESS"`
	FileStorePath string `json:"file_storage_path" yaml:"FILE_STORAGE_PATH"`
	DataBaseDSN   string `json:"database_dsn" yaml:"DATABASE_DSN"`
	EnableHTTPS   bool   `json:"enable_https" yaml:"ENABLE_HTTPS"`
	TrustedSubnet string `json:"trusted_subnet" yaml:"TRUSTED_SUBNET"`
	ServerRPC     string `json:"server_address_grpc" yaml:"SERVER_GRPC"`
}
type ConfigFromENV struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FileStorePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN   string `env:"DATABASE_DSN"`
	EnableHTTPS   string `env:"ENABLE_HTTPS"`
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
	ServerRPC     string `env:"SERVER_GRPC"`
}
type ServerConfig interface {
	GetConfig() (config *ServiceShortURLConfig, err error)
}

func GetConfig() (config *ServiceShortURLConfig, err error) {

	var ServerAddress, BaseURL, FilePath, configFileName, TrustedSubnet, DataBaseDSN, ServerRPC string
	EnableHTTPS := false
	flag.StringVarP(&configFileName, "config", "c", "", "path to the configuration file")
	flag.StringVarP(&ServerAddress, "server", "a", "", "Server Adders HTTP")
	flag.StringVarP(&ServerRPC, "server_rpc", "g", "", "Server Adders GRPC")
	flag.StringVarP(&BaseURL, "base_url", "b", "", "Base URL")
	flag.StringVarP(&TrustedSubnet, "trusted_subnet", "t", "", "Base URL")
	flag.StringVarP(&FilePath, "file_path", "f", "", "Config file path")
	flag.StringVarP(&DataBaseDSN, "dsn", "d", "", "Set driver DSN ")
	flag.BoolVarP(&EnableHTTPS, "tls", "s", false, "In HTTP server is Enable TLS")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	log.Infof("Loading configuration at '%s'", configFileName)
	configFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Errorf("can't open the config file: %s", err)

	}
	// Default values
	config = &ServiceShortURLConfig{
		Port:          8080,
		ServerHost:    "localhost",
		BaseURL:       "http://localhost:8080",
		ServerAddress: "localhost:8080",
		FileStorePath: "",
		EnableHTTPS:   false,
		DataBaseDSN:   "",
		TrustedSubnet: "",
		ServerRPC:     "localhost:3200",
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Errorf("YAML can't read the config file: %s", err)
	}

	var cfgenv ConfigFromENV
	err = env.Parse(&cfgenv)
	if err != nil {
		log.Errorf("can't start the listening thread: %s", err)
	} else {
		if cfgenv.ServerAddress != "" {
			config.ServerAddress = cfgenv.ServerAddress
		}
		if cfgenv.ServerRPC != "" {
			config.ServerRPC = cfgenv.ServerRPC
		}
		if cfgenv.BaseURL != "" {
			config.BaseURL = cfgenv.BaseURL
		}
		if cfgenv.FileStorePath != "" {
			config.FileStorePath = cfgenv.FileStorePath
		}
		if cfgenv.DataBaseDSN != "" {
			config.DataBaseDSN = cfgenv.DataBaseDSN
		}
		if cfgenv.EnableHTTPS != "" {
			config.EnableHTTPS = true
		}
		if cfgenv.TrustedSubnet != "" {
			config.TrustedSubnet = cfgenv.TrustedSubnet
		}
	}

	if ServerAddress != "" {
		config.ServerAddress = ServerAddress
	}
	if ServerRPC != "" {
		config.ServerRPC = ServerRPC
	}
	if BaseURL != "" {
		config.BaseURL = BaseURL
	}
	if FilePath != "" {
		config.FileStorePath = FilePath
	}
	if DataBaseDSN != "" {
		config.DataBaseDSN = DataBaseDSN
	}
	if TrustedSubnet != "" {
		config.TrustedSubnet = TrustedSubnet
	}
	if EnableHTTPS {
		config.EnableHTTPS = EnableHTTPS
	}
	//***postgres:5432/praktikum?sslmode=disable

	if configFileName != "" {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			log.Errorf("JSON can't read the config file: %s", err)
		}
	}
	log.Info(config.DataBaseDSN)
	log.Info("Configuration loaded")
	return
}
