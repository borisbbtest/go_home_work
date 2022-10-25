package app

import (
	"net"

	"github.com/borisbbtest/go_home_work/internal/config"
	handlersgrpc "github.com/borisbbtest/go_home_work/internal/handlers/handlers_grpc"
	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"google.golang.org/grpc"
)

type serviceRPCShortURL struct {
	wrapp handlersgrpc.WrapperHandlerRPC
}

func NewRPC(cfg *config.ServiceShortURLConfig) *serviceRPCShortURL {
	return &serviceRPCShortURL{
		wrapp: handlersgrpc.WrapperHandlerRPC{
			ServerConf: cfg,
		},
	}
}

func (hook *serviceRPCShortURL) Start() (err error) {

	hook.wrapp.Storage, err = storage.NewPostgreSQLStorage(hook.wrapp.ServerConf.DataBaseDSN)
	if err != nil {
		hook.wrapp.Storage, err = storage.NewFileStorage(hook.wrapp.ServerConf.FileStorePath)
		if err != nil {
			log.Error(err)
		}
	}

	listen, err := net.Listen("tcp", hook.wrapp.ServerConf.ServerRPC)
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	shortrpc.RegisterShortURLServer(s, &hook.wrapp)

	log.Info("Server gRPC is running ")

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	log.Info("End GRPC")
	defer s.Stop()
	defer listen.Close()
	defer hook.wrapp.Storage.Close()
	return
}
