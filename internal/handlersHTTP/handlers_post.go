package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type WrapperHandler struct {
	ServerConf *config.ServiceShortURLConfig
	Storage    storage.Storage
	UserID     string
}

// PostHandler делает полезные дела ))
func (hook *WrapperHandler) PostHandler(w http.ResponseWriter, r *http.Request) {

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}

	log.Info("PostHandler ", string(bytes))

	hashcode, _ := storage.ParserDataURL(string(bytes))

	hashcode.UserID = hook.UserID

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	gl, err := hook.Storage.Put(hashcode.ShortPath, hashcode)

	if err != nil {
		log.Error("Put error ", err)
	}

	if len(gl) > 1 {
		hashcode.ShortPath = gl
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	resp := fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath)

	fmt.Fprint(w, resp)

	log.Println("Post handler")
	defer r.Body.Close()
}

// PostJSONHandler и тут тоже
func (hook *WrapperHandler) PostJSONHandler(w http.ResponseWriter, r *http.Request) {

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	log.Info("PostJSONHandler")
	defer r.Body.Close()

	var m model.RequestAddDBURL
	if err := json.Unmarshal(bytes, &m); err != nil {
		log.Errorf("body error: %v", string(bytes))
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "request body is not valid json", 400)
		return
	}

	hashcode, err := storage.ParserDataURL(m.ReqNewURL)
	if err != nil {
		http.Error(w, "request body is not valid URL", 400)
		return
	}

	hashcode.UserID = hook.UserID

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ShortPath, err := hook.Storage.Put(hashcode.ShortPath, hashcode)
	if err != nil {
		log.Error("Put error ", err)
	}

	// проверяем что получили хеш сокращенного url
	if ShortPath != "" {
		hashcode.ShortPath = ShortPath
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := model.ResponseURLShort{
		ResNewURL: fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath),
	}

	json.NewEncoder(w).Encode(resp)

	log.Println("Post handler")
}

// PostJSONHandlerBatch  осмыслено  собирает много чего
func (hook *WrapperHandler) PostJSONHandlerBatch(w http.ResponseWriter, r *http.Request) {

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	log.Info("PostJSONHandler")
	defer r.Body.Close()

	var m []model.RequestBatch
	if err := json.Unmarshal(bytes, &m); err != nil {
		log.Errorf("body error: %v", string(bytes))
		log.Errorf("error decoding message: %v", err)
		http.Error(w, "request body is not valid json", 400)
		return
	}

	res1, res2 := storage.ParserDataURLBatch(&m, hook.ServerConf.BaseURL, hook.UserID)
	if err != nil {
		http.Error(w, "request body is not valid URL", 400)
		return
	}
	hook.Storage.PutBatch(hook.UserID, res1)

	// resp := model.ResponseURLShort{
	// 	ResNewURL: fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath),
	// }

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res2)

	log.Println("Post handler")
}
