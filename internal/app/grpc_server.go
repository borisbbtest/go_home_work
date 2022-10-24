package app

import (
	"net"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/handlers"
	"google.golang.org/grpc"
)

type service_RPC_ShortURL struct {
	wrapp handlers.WrapperHandler
}

func NewRPC(cfg *config.ServiceShortURLConfig) *serviceHttpShortURL {
	return &serviceHttpShortURL{
		wrapp: handlers.WrapperHandler{
			ServerConf: cfg,
		},
	}
}

func (hook *service_RPC_ShortURL) Start() (err error) {

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	//pb.RegisterUsersServer(s, &UsersServer{})

	log.Info("Server gRPC is running ")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

}
