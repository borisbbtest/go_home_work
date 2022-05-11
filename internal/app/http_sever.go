package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
		},
	}
}

func (hook *serviceShortURL) Start() error {

	// Launch the listening thread
	log.Println("Initializing HTTP server")
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", hook.wrapp.GetHandler)
	r.Post("/", hook.wrapp.PostHandler)
	r.Post("/api/shorten", hook.wrapp.PostJSONHandler)
	r.Get("/", hook.wrapp.GetHandler)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	hook.wrapp.FileServer(r, "/form", filesDir)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(hook.conf.Port)),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("can't start the listening thread: %s", err)
	}

	log.Info("Exiting")
	return nil
}
