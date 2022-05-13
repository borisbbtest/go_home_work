package storage

import (
	"fmt"
	"net"
	"net/url"

	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/sirupsen/logrus"
)

type StorageURL struct {
	Port      string `json:"Port"`
	URL       string `json:"URL"`
	Path      string `json:"Path"`
	ShortPath string `json:"ShortPath"`
}

type StoreDB struct {
	DBLocal map[string]StorageURL
}
type IStorageURL interface {
	StoreDBinMemory(str string) (res string, err error)
	GetShortURLfromDBinMemory(urlshort string) (resdirect string)
}

var log = logrus.WithField("context", "service_short_url")

func (store *StoreDB) StoreDBinMemory(str string) (res StorageURL, err error) {

	url, err := url.ParseRequestURI(str)
	if err != nil {
		log.Info(err.Error())
		return res, err
	}
	hesh := tools.GenerateShortLink(str)
	address := net.ParseIP(url.Host)
	log.Println("url-info", "host", address)
	res.Path = url.Path
	res.Port = url.Port()
	res.URL = str
	res.ShortPath = hesh
	store.DBLocal[hesh] = res
	return
}

func (store *StoreDB) GetShortURLfromDBinMemory(urlshort string) (resdirect string) {
	if _, ok := store.DBLocal[urlshort]; ok {
		return fmt.Sprintf("%s", store.DBLocal[urlshort])
	}
	return
}
