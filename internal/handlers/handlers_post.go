package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type WrapperHandler struct {
	ServerConf *config.ServiceShortURLConfig
	Storage    storage.Storage
}
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func (hook *WrapperHandler) GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

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
	resp := fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath)

	tmp, _ := tools.GetID()
	v, _ := tools.AddCookie(w, r, "ShortURL", fmt.Sprintf("%x", tmp), 30*time.Minute)
	hashcode.UserID = v

	w.WriteHeader(http.StatusCreated)
	if v, err := hook.Storage.Put(hashcode.ShortPath, hashcode); err == nil {
		hashcode.ShortPath = v
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, resp)

	log.Println("Post handler")
	defer r.Body.Close()
}

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

	tmp, _ := tools.GetID()
	v, _ := tools.AddCookie(w, r, "ShortURL", fmt.Sprintf("%x", tmp), 30*time.Minute)
	hashcode.UserID = v

	if v, err := hook.Storage.Put(hashcode.ShortPath, hashcode); err == nil {
		hashcode.ShortPath = v
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := model.ResponseURLShort{
		ResNewURL: fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)

	log.Println("Post handler")
	defer r.Body.Close()
}

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

	tmp, _ := tools.GetID()
	v, _ := tools.AddCookie(w, r, "ShortURL", fmt.Sprintf("%x", tmp), 30*time.Minute)

	res1, res2 := storage.ParserDataURLBatch(&m, hook.ServerConf.BaseURL, v)
	if err != nil {
		http.Error(w, "request body is not valid URL", 400)
		return
	}
	hook.Storage.PutBatch(v, res1)

	// resp := model.ResponseURLShort{
	// 	ResNewURL: fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath),
	// }

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res2)

	log.Println("Post handler")
	defer r.Body.Close()
}
