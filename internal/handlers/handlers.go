package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
)

type WrapperHandler struct {
	urlStore storage.StorageURL
}

var log = logrus.WithField("context", "service_short_url")

func (hook *WrapperHandler) GetHandler(w http.ResponseWriter, r *http.Request) {

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
	fmt.Printf(m)
}

func (hook *WrapperHandler) PostHandler(w http.ResponseWriter, r *http.Request) {

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
	fmt.Printf(m)
}
