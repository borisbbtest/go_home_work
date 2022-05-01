package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/handlers"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type service_short_url struct {
	conf  Service_short_urlConfig
	wrapp handlers.WrapperHandler
}

func New(cfg *Service_short_urlConfig) *service_short_url {
	return &service_short_url{
		conf: *cfg,
		wrapp: handlers.WrapperHandler{
			UrlStore: storage.StoreDB{
				DBLocal: make(map[string]storage.StorageURL),
			},
			ServerConf: cfg,
		},
	}
}

func (hook *service_short_url) Start() error {

	// Launch the listening thread
	log.Println("Initializing HTTP server 1")
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", hook.wrapp.GetHandler)
	r.Post("/", hook.wrapp.PostHandler)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	hook.wrapp.FileServer(r, "/", filesDir)

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
