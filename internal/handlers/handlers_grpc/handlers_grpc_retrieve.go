package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

type WrapperHandlerRPC struct {
	ServerConf *config.ServiceShortURLConfig
	Storage    storage.Storage
	UserID     string
	shortrpc.UnimplementedShortURLServer
}

func (hook *WrapperHandlerRPC) Retrieve(ctx context.Context, in *shortrpc.RetrieveRequest) (*shortrpc.RetrieveResponse, error) {
	res := shortrpc.RetrieveResponse{
		Status: "ok",
	}
	value, err := hook.Storage.Get(in.GetShortUrlId())

	if err != nil {
		res.Status = err.Error()
	}
	res.RedirectUrl = value.URL

	return &res, nil
}
