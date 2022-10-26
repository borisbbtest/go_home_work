package handlersgrpc

import (
	"context"
	"strconv"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

func (hook *WrapperHandlerRPC) CreateBatch(ctx context.Context, in *shortrpc.CreateBatchRequest) (*shortrpc.CreateBatchResponse, error) {

	res := shortrpc.CreateBatchResponse{
		Status: "OK",
	}
	urls := []model.RequestBatch{}

	for _, v := range in.GetUrls() {
		str := strconv.FormatInt(int64(v.CorrelationId), 10)
		urls = append(urls, model.RequestBatch{
			CorrelationID: str,
			OriginalURL:   v.GetOriginalUrl(),
		})
	}

	res1, res2 := storage.ParserDataURLBatch(&urls, hook.ServerConf.BaseURL, in.GetUserId())

	err := hook.Storage.PutBatch(in.GetUserId(), res1)

	if err != nil {
		res.Status = err.Error()
		return &res, status.Error(codes.FailedPrecondition, err.Error())
	}

	for _, v := range res2 {
		intVar, _ := strconv.Atoi(v.CorrelationID)
		res.Urls = append(res.Urls, &shortrpc.CreateBatchResponse_URL{CorrelationId: int32(intVar), ShortUrl: v.ShortURL})
	}

	return &res, nil
}
