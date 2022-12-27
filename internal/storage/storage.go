package storage

import (
	"fmt"
	"net"
	"net/url"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/sirupsen/logrus"
)

const (
	RecordNotFound = 2
)

type StoreDBLocal struct {
	DBLocal  map[string]model.DataURL
	ListUser map[string][]string
}
type Storage interface {
	Put(k string, v model.DataURL) (string, error)
	Get(k string) (model.DataURL, error)
	GetAll(k string, dom string) ([]model.ResponseURL, error)
	GetStats() (model.ResponseStats, error)
	PutBatch(k string, v []model.DataURL) error
	DeletedURLBatch(k string, v []model.DataURL) error
	Close()
}

var log = logrus.WithField("context", "service_short_url")

// Обертка к запросу
func ParserDataURL(str string) (res model.DataURL, err error) {

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

// Обертка к запросу
func ParserDataURLBatch(str *[]model.RequestBatch, basedurl string, userid string) (res []model.DataURL, res2 []model.ResponseBatch) {
	res = []model.DataURL{}
	res2 = []model.ResponseBatch{}
	for _, element := range *str {
		if k, err := ParserDataURL(element.OriginalURL); err == nil {
			k.CorrelationID = element.CorrelationID
			k.UserID = userid
			res = append(res, k)
			res2 = append(res2, model.ResponseBatch{
				CorrelationID: element.CorrelationID,
				ShortURL:      fmt.Sprintf("%s/%s", basedurl, k.ShortPath),
			})
		}
	}
	return
}
