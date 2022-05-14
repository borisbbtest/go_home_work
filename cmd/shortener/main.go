package main

import (
	"os"

	goflag "flag"

	"github.com/borisbbtest/go_home_work/internal/app"
	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var log = logrus.WithField("context", "main")

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	configFileName := flag.String("config", "./config.yml", "path to the configuration file")
	var ServerAddress, BaseURL, FilePath string
	flag.StringVarP(&ServerAddress, "server", "a", "", "Server Adders")
	flag.StringVarP(&BaseURL, "base_url", "b", "", "Base URL")
	flag.StringVarP(&FilePath, "file_path", "f", "", "Config file path")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	cfg, err := config.GetConfig(*configFileName)
	if err != nil {
		cfg = &config.ServiceShortURLConfig{
			Port:          8080,
			ServerHost:    "localhost",
			BaseURL:       "http://localhost:8080",
			ServerAddress: "localhost:8080",
			FileStorePath: "",
		}
	}

	//  получаем переменные среды
	var cfgenv config.ConfigFromENV
	e := env.Parse(&cfgenv)
	if e != nil {
		log.Errorf("can't start the listening thread: %s", e)
	} else {
		if cfgenv.ServerAddress != "" {
			cfg.ServerAddress = cfgenv.ServerAddress
		}
		if cfgenv.BaseURL != "" {
			cfg.BaseURL = cfgenv.BaseURL
		}
		if cfgenv.FileStorePath != "" {
			cfg.FileStorePath = cfgenv.FileStorePath
		}
	}

	if ServerAddress != "" {
		cfg.ServerAddress = ServerAddress
	}
	if BaseURL != "" {
		cfg.BaseURL = BaseURL
	}
	if FilePath != "" {
		cfg.FileStorePath = FilePath
	}

	err = app.New(cfg).Start()
	if err != nil {
		log.Fatal(err)
	}
}
