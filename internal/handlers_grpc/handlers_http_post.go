package handlersgrpc

import (
	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type WrapperHandler struct {
	ServerConf *config.ServiceShortURLConfig
	Storage    storage.Storage
	UserID     string
}
