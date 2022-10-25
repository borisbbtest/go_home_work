package app

import (
	"net"

	"github.com/borisbbtest/go_home_work/internal/config"
	handlersgrpc "github.com/borisbbtest/go_home_work/internal/handlers_grpc"
	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"google.golang.org/grpc"
)

type service_RPC_ShortURL struct {
	wrapp handlersgrpc.WrapperHandler
}

func NewRPC(cfg *config.ServiceShortURLConfig) *service_RPC_ShortURL {
	return &service_RPC_ShortURL{
		wrapp: handlersgrpc.WrapperHandler{
			ServerConf: cfg,
		},
	}
}

func (hook *service_RPC_ShortURL) Start() (err error) {

	hook.wrapp.Storage, err = storage.NewPostgreSQLStorage(hook.wrapp.ServerConf.DataBaseDSN)
	if err != nil {
		hook.wrapp.Storage, err = storage.NewFileStorage(hook.wrapp.ServerConf.FileStorePath)
		if err != nil {
			log.Error(err)
		}
	}

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	shortrpc.RegisterShortURLServer(s, &handlersgrpc.WrapperHandler{})

	log.Info("Server gRPC is running ")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	defer hook.wrapp.Storage.Close()
	return
}
