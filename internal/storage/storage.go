package storage

import (
	"fmt"
	"net"
	"net/url"

	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/sirupsen/logrus"
)

type StorageURL struct {
	Port string
	Url  string
	Path string
}

type StoreDB struct {
	DBLocal map[string]StorageURL
}
type IStorageURL interface {
	setURLforRedirect() (err error)
	getURLforRedirect(err error, urlshort string)
}

var log = logrus.WithField("context", "service_short_url")

func (store *StoreDB) PostURLforRedirect(str string) (res string, err error) {

	url, err := url.ParseRequestURI(str)
	if err != nil {
		log.Info(err.Error())
		return res, err
	}
	res = tools.GenerateShortLink(str)
	address := net.ParseIP(url.Host)
	log.Println("url-info", "host", address)
	var dataStore StorageURL
	dataStore.Path = url.Path
	dataStore.Port = url.Port()
	dataStore.Url = str
	store.DBLocal[res] = dataStore
	return
}

func (store *StoreDB) GetURLforRedirect(urlshort string) (resdirect string) {
	if _, ok := store.DBLocal["foo"]; ok {
		return fmt.Sprintf("%s", store.DBLocal[urlshort])
	}
	return
}
