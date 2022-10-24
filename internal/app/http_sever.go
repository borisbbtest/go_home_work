package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/borisbbtest/go_home_work/internal/config"
	handlershttp "github.com/borisbbtest/go_home_work/internal/handlers_http"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type serviceHttpShortURL struct {
	wrapp handlershttp.WrapperHandler
}

// Структура так так
func NewHTTP(cfg *config.ServiceShortURLConfig) *serviceHttpShortURL {
	return &serviceHttpShortURL{
		wrapp: handlershttp.WrapperHandler{
			ServerConf: cfg,
		},
	}
}

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func printIntro() {
	log.Info("Build version: ", buildVersion)
	log.Info("Build date: ", buildDate)
	log.Info("Build commit: ", buildCommit)
}
func (hook *serviceHttpShortURL) Start() (err error) {

	printIntro()
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
	//	defer hook.wrapp.Storage.Close()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(hook.wrapp.GzipHandle)
	r.Use(hook.wrapp.MidSetCookie)
	//r.Use(middleware.Compress(5, "gzip"))
	r.Use(middleware.Recoverer)
	//yes
	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	r.Get("/api/user/urls", hook.wrapp.GetHandlerCooke)
	r.Get("/", hook.wrapp.GetHandler)
	r.Get("/ping", hook.wrapp.GetHandlerPing)
	r.Get("/{id}", hook.wrapp.GetHandler)
	r.Post("/", hook.wrapp.PostHandler)
	r.Post("/api/shorten", hook.wrapp.PostJSONHandler)
	r.Post("/api/shorten/batch", hook.wrapp.PostJSONHandlerBatch)
	r.Delete("/api/user/urls", hook.wrapp.DeleteURLHandlers)
	r.Get("/api/internal/stats", hook.wrapp.GetHandlerStats)

	server := &http.Server{
		Addr:         hook.wrapp.ServerConf.ServerAddress,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		for {
			s := <-sigint
			switch s {
			case syscall.SIGINT:
				if err := server.Shutdown(context.Background()); err != nil {
					// ошибки закрытия Listener
					log.Printf("HTTP server Shutdown SIGINT:  %v", err)
				}
				log.Info("bz -SIGINT")
				close(idleConnsClosed)
			case syscall.SIGTERM:
				if err := server.Shutdown(context.Background()); err != nil {
					// ошибки закрытия Listener
					log.Printf("HTTP server Shutdown SIGTERM: %v", err)
				}
				log.Info("bz - SIGTERM")
				close(idleConnsClosed)
			case syscall.SIGQUIT:
				if err := server.Shutdown(context.Background()); err != nil {
					// ошибки закрытия Listener
					log.Printf("HTTP server Shutdown SIGQUIT: %v", err)
				}
				log.Info("bz - SIGQUIT")
				close(idleConnsClosed)
			default:
				fmt.Println("Unknown signal.")
			}
		}
	}()

	defer server.Close()
	if hook.wrapp.ServerConf.EnableHTTPS {

		cert, key, err := tools.CertGeg()
		if err != nil {
			return fmt.Errorf("BZ Certificate and key wasn't generation: %s", err)
		}

		tools.WriteCertFile("cert.pem", cert)
		tools.WriteCertFile("key.pem", key)
		err = server.ListenAndServeTLS("cert.pem", "key.pem")

		if err != http.ErrServerClosed {
			return fmt.Errorf("BZ can't start the listening thread: %s", err)
		}

	} else {
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {

			return fmt.Errorf("BZ can't start the listening thread: %s", err)
		}

	}
	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы
	hook.wrapp.Storage.Close()

	log.Info("Server Shutdown gracefully")

	log.Info("Exiting")
	return nil
}
