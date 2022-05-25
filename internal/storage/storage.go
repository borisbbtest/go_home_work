package storage

import (
	"net"
	"net/url"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/sirupsen/logrus"
)

type DataURL struct {
	Port      string `json:"Port"`
	URL       string `json:"URL"`
	Path      string `json:"Path"`
	ShortPath string `json:"ShortPath"`
	UserID    string `json:"UserID"`
}
type StoreDBLocal struct {
	DBLocal  map[string]DataURL
	ListUser map[string][]string
}
type Storage interface {
	Put(k string, v DataURL) error
	Get(k string) (DataURL, error)
	GetAll(k string, dom string) ([]model.ResponseURL, error)
	Close()
}

var log = logrus.WithField("context", "service_short_url")

func ParserDataURL(str string) (res DataURL, err error) {

	url, err := url.ParseRequestURI(str)
	if err != nil {
		log.Info(err.Error())
		return res, err
	}
	hash := tools.GenerateShortLink(str)
	address := net.ParseIP(url.Host)
	log.Println("url-info", "host", address)
	res.Path = url.Path
	res.Port = url.Port()
	res.URL = str
	res.ShortPath = hash
	return
}
