package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
)

func (hook *WrapperHandlerRPC) GetUserURLs(ctx context.Context, in *shortrpc.GetUserURLsRequest) (*shortrpc.GetUserURLsResponse, error) {
	res := shortrpc.GetUserURLsResponse{
		Status: "Ok",
	}

	responseShortURL, err := hook.Storage.GetAll(in.UserId, hook.ServerConf.BaseURL)
	log.Info("Size  - ", len(responseShortURL))
	if err != nil {
		res.Status = err.Error()
	}
	for index, value := range responseShortURL {
		log.Info(value.OriginalURL)
		res.Urls[index].OriginalUrl = value.OriginalURL
		res.Urls[index].ShortUrl = value.ShortURL
	}

	return &res, nil
}
