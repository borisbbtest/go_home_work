package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
)

func (hook *WrapperHandlerRPC) DeleteBatch(ctx context.Context, in *shortrpc.DeleteBatchRequest) (*shortrpc.DeleteBatchResponse, error) {

	res := shortrpc.DeleteBatchResponse{
		Status: "Accept",
	}
	go func() {
		buff := make([]model.DataURL, 0, len(in.GetUrls()))
		for _, str := range in.GetUrls() {
			buff = append(buff, model.DataURL{ShortPath: str,
				StatusActive: 2})
		}
		hook.Storage.DeletedURLBatch(in.GetUserId(), buff)
	}()

	return &res, nil
}
