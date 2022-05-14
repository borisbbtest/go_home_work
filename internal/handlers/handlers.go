package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type WrapperHandler struct {
	URLStore   storage.StoreDB
	ServerConf *config.ServiceShortURLConfig
	FielDB     *storage.InitStoreDBinFile
}
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

var log = logrus.WithField("context", "service_short_url")

func (hook *WrapperHandler) GetHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	// for k, v := range hook.UrlStore.DBLocal {
	// 	fmt.Printf("key[%s] value[%s]\n", k, v.Url)
	// }
	log.Info("ID Go to", id)
	value, status := hook.URLStore.DBLocal[id]
	if status {
		url := value.URL
		w.Header().Set("Location", url)
		w.WriteHeader(307)
		//log.Printf("Get handler")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
		//log.Printf("key not found")
	}
	fmt.Println(id)
	defer r.Body.Close()
	//log.Printf("Get handler")
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
	hashcode, _ := hook.URLStore.StoreDBinMemory(string(bytes))
	resp := fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath)

	if hook.FielDB != nil {
		if hook.FielDB.WriteURL != nil {
			hook.FielDB.WriteURL.WriteEvent(&hashcode)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)
	fmt.Fprint(w, resp)

	log.Println("Post handler")
	defer r.Body.Close()
}

func (hook *WrapperHandler) PostJSONHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
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

	hashcode, err := hook.URLStore.StoreDBinMemory(m.ReqNewURL)
	if err != nil {
		http.Error(w, "request body is not valid URL", 400)
		return
	}

	if hook.FielDB != nil {
		if hook.FielDB.WriteURL != nil {
			hook.FielDB.WriteURL.WriteEvent(&hashcode)
		}
	}
	resp := model.ResponseURLShort{
		ResNewURL: fmt.Sprintf("%s/%s", hook.ServerConf.BaseURL, hashcode.ShortPath),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	log.Println("Post handler")
	defer r.Body.Close()
}

func (hook *WrapperHandler) FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"
	log.Println(path)
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
