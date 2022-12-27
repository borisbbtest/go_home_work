package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

func (hook *WrapperHandlerRPC) GetUserURLs(ctx context.Context, in *shortrpc.GetUserURLsRequest) (*shortrpc.GetUserURLsResponse, error) {
	res := shortrpc.GetUserURLsResponse{
		Status: "Ok",
	}

	responseShortURL, err := hook.Storage.GetAll(in.UserId, hook.ServerConf.BaseURL)
	log.Info("Size  - ", len(responseShortURL))
	if err != nil {
		res.Status = err.Error()
		return &res, status.Error(codes.FailedPrecondition, err.Error())
	}
	for _, value := range responseShortURL {
		x := shortrpc.GetUserURLsResponse_URL{
			ShortUrl:    value.ShortURL,
			OriginalUrl: value.OriginalURL,
		}
		res.Urls = append(res.Urls, &x)

	}

	return &res, nil
}
