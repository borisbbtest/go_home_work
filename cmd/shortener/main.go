package main

import (
	"os"

	_ "net/http/pprof"

	"github.com/borisbbtest/go_home_work/internal/app"
	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "main")

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	cfg, err := config.GetConfig()
	if err != nil {
		cfg = &config.ServiceShortURLConfig{
			Port:          8080,
			ServerHost:    "localhost",
			BaseURL:       "http://localhost:8080",
			ServerAddress: "localhost:8080",
			FileStorePath: "",
		}
	}
	err = app.New(cfg).Start()
	if err != nil {
		log.Fatal(err)
	}
}
