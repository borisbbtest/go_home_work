package app

import (
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
	wrapp handlers.WrapperHandler
}

func New(cfg *config.ServiceShortURLConfig) *serviceShortURL {
	return &serviceShortURL{
		wrapp: handlers.WrapperHandler{
			ServerConf: cfg,
		},
	}
}

func (hook *serviceShortURL) Start() (err error) {

	// Launch the listening thread
	log.Println("Initializing HTTP server")
	r := chi.NewRouter()

	hook.wrapp.Storage, err = storage.NewPostgreSQLStorage(hook.wrapp.ServerConf.DataBaseDSN)
	if err != nil {
		hook.wrapp.Storage, err = storage.NewFileStorage(hook.wrapp.ServerConf.FileStorePath)
		if err != nil {
			log.Error(err)
		}
	}
	defer hook.wrapp.Storage.Close()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(hook.wrapp.GzipHandle)
	r.Use(hook.wrapp.MidSetCookie)
	//r.Use(middleware.Compress(5, "gzip"))
	r.Use(middleware.Recoverer)

	r.Get("/{id}", hook.wrapp.GetHandler)
	r.Post("/", hook.wrapp.PostHandler)
	r.Post("/api/shorten", hook.wrapp.PostJSONHandler)
	r.Post("/api/shorten/batch", hook.wrapp.PostJSONHandlerBatch)
	r.Get("/api/user/urls", hook.wrapp.GetHandlerCooke)
	r.Get("/", hook.wrapp.GetHandler)
	r.Get("/ping", hook.wrapp.GetHandlerPing)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	hook.wrapp.FileServer(r, "/form", filesDir)

	server := &http.Server{
		Addr:         hook.wrapp.ServerConf.ServerAddress,
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
