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

func NewRPC(cfg *config.ServiceShortURLConfig, st storage.Storage) *serviceRPCShortURL {
	return &serviceRPCShortURL{
		wrapp: handlersgrpc.WrapperHandlerRPC{
			ServerConf: cfg,
			Storage:    st,
		},
	}
}

func (hook *serviceRPCShortURL) Start() (err error) {

	listen, err := net.Listen("tcp", hook.wrapp.ServerConf.ServerRPC)
	if err != nil {
		log.Fatal(err)
		return
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	shortrpc.RegisterShortURLServer(s, &hook.wrapp)

	log.Info("Server gRPC is running ")

	// получаем запрос gRPC
	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info("End GRPC")
	defer s.Stop()
	defer listen.Close()
	defer hook.wrapp.Storage.Close()
	return
}
