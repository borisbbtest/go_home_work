package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type WrapperHandler struct {
	URLStore   storage.StoreDB
	ServerConf *config.ServiceShortURLConfig
}

var log = logrus.WithField("context", "service_short_url")

func (hook *WrapperHandler) GetHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	// for k, v := range hook.UrlStore.DBLocal {
	// 	fmt.Printf("key[%s] value[%s]\n", k, v.Url)
	// }

	value, status := hook.URLStore.DBLocal[id]
	if status {
		url := value.URL
		w.Header().Set("Location", url)
		w.WriteHeader(307)
		log.Printf("Get handler")
	} else {
		w.WriteHeader(201)
		fmt.Fprintln(w, "OK")
		log.Printf("key not found")
	}

	fmt.Println(id)
	defer r.Body.Close()

	log.Printf("Get handler")
}

func (hook *WrapperHandler) PostHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	hashcode, _ := hook.URLStore.PostURLforRedirect(string(bytes))
	resp := fmt.Sprintf("http://%s:%d/%s", hook.ServerConf.ServerHost, hook.ServerConf.Port, hashcode)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)
	fmt.Fprintln(w, resp)

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
