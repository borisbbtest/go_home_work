package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type WrapperHandler struct {
	ServerConf *config.ServiceShortURLConfig
	Storage    storage.Storage
	UserID     string
	shortrpc.UnimplementedShortURLServer
}

func (hool *WrapperHandler) Retrieve(ctx context.Context, in *shortrpc.RetrieveRequest) (*shortrpc.RetrieveResponse, error) {
	var res shortrpc.RetrieveResponse
	res = shortrpc.RetrieveResponse{
		Status:      "1",
		RedirectUrl: "sss",
	}
	return &res, nil
}
