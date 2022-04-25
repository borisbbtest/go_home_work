package main

import (
	"flag"
	"os"

	. "github.com/borisbbtest/go_home_work/internal/app"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "main")

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	configFileName := flag.String("config", "./config.yml", "path to the configuration file")
	flag.Parse()

	cfg, err := ConfigFromFile(*configFileName)
	if err != nil {
		log.Fatal(err)
	}

	err = New(cfg).Start()
	if err != nil {
		log.Fatal(err)
	}
}
