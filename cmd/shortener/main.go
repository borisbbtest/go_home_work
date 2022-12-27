package main

import (
	"os"

	"github.com/borisbbtest/go_home_work/internal/app"
	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "main")

// Использую в  vsc
// main "go.formatTool": "gofmt"
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
	ap, err := app.Init(cfg)
	if err != nil {
		log.Error(err)
	}
	ap.Start()
}
