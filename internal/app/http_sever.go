package app

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/handlers"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type serviceShortURL struct {
	conf  config.ServiceShortURLConfig
	wrapp handlers.WrapperHandler
}

func New(cfg *config.ServiceShortURLConfig) *serviceShortURL {
	return &serviceShortURL{
		conf: *cfg,
		wrapp: handlers.WrapperHandler{
			URLStore: storage.StoreDB{
				DBLocal: make(map[string]storage.StorageURL),
			},
			ServerConf: cfg,
			FielDB:     &storage.InitStoreDBinFile{},
		},
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
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

func (hook *serviceShortURL) Start() (err error) {

	// Launch the listening thread
	log.Println("Initializing HTTP server")
	r := chi.NewRouter()

	if hook.conf.FileStorePath != "" {
		hook.wrapp.FielDB.ReadURL, err = storage.NewConsumer(hook.conf.FileStorePath)
		if err != nil {
			log.Fatal(err)
		}
		defer hook.wrapp.FielDB.ReadURL.Close()

		hook.wrapp.FielDB.WriteURL, err = storage.NewProducer(hook.conf.FileStorePath)
		if err != nil {
			log.Fatal(err)
		}
		defer hook.wrapp.FielDB.WriteURL.Close()

		scanner := bufio.NewScanner(hook.wrapp.FielDB.ReadURL.GetFile())
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			var m storage.StorageURL
			if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
				log.Errorf("body error: %v", string(scanner.Bytes()))
				log.Errorf("error decoding message: %v", err)
			}
			hook.wrapp.URLStore.DBLocal[m.ShortPath] = m
		}

	}
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(gzipHandle)
	r.Use(middleware.Compress(5, "gzip"))
	r.Use(middleware.Recoverer)

	r.Get("/{id}", hook.wrapp.GetHandler)
	r.Post("/", hook.wrapp.PostHandler)
	r.Post("/api/shorten", hook.wrapp.PostJSONHandler)
	r.Get("/", hook.wrapp.GetHandler)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	hook.wrapp.FileServer(r, "/form", filesDir)

	server := &http.Server{
		Addr:         hook.conf.ServerAddress,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("can't start the listening thread: %s", err)
	}

	log.Info("Exiting")
	return nil
}
