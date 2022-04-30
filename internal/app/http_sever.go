package app

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/handlers"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type service_short_url struct {
	conf  Service_short_urlConfig
	wrapp handlers.WrapperHandler
}

func New(cfg *Service_short_urlConfig) *service_short_url {

	return &service_short_url{
		conf: *cfg,
	}
}

func (hook *service_short_url) Start() error {

	// Launch the listening thread
	log.Println("Initializing HTTP server")
	http.HandleFunc("/", hook.wrapp.GetHandler)
	http.HandleFunc("/", hook.wrapp.PostHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(hook.conf.Port), nil)
	if err != nil {
		return fmt.Errorf("can't start the listening thread: %s", err)
	}

	log.Info("Exiting")
	return nil
}
