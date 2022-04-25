package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

func (hook *Service_short_url) MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hook.PostHandler(w, r)
	case http.MethodGet:
		hook.GetHandler(w, r)
	default:
		http.Error(w, "unsupported HTTP method only post send", 400)
	}
}

func (hook *Service_short_url) GetHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(bytes))

	defer r.Body.Close()
	var m string
	if err := json.Unmarshal(bytes, &m); err != nil {
		log.Errorf("body error: %v", string(bytes))
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "request body is not valid json", 400)
		return
	}

	hook.ChannelGet <- &m
	fmt.Printf(m)

}

func (hook *Service_short_url) PostHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(bytes))

	defer r.Body.Close()
	var m string
	if err := json.Unmarshal(bytes, &m); err != nil {
		log.Errorf("body error: %v", string(bytes))
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "request body is not valid json", 400)
		return
	}

	hook.ChannelPost <- &m
	fmt.Printf(m)

}
