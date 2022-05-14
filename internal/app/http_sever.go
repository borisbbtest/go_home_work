package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	//r.Use(hook.wrapp.GzipHandle)
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
