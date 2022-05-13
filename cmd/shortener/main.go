package main

import (
	"flag"
	"os"

	"github.com/borisbbtest/go_home_work/internal/app"
	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "main")

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	configFileName := flag.String("config", "./config.yml", "path to the configuration file")
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
		if cfgenv.ServerAddress != "" {
			cfg.BaseURL = cfgenv.BaseURL
		}
		if cfgenv.FileStorePath != "" {
			cfg.FileStorePath = cfgenv.FileStorePath
		}
	}

	err = app.New(cfg).Start()
	if err != nil {
		log.Fatal(err)
	}
}
