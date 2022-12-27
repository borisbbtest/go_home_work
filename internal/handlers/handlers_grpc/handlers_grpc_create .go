package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

func (hook *WrapperHandlerRPC) Create(ctx context.Context, in *shortrpc.CreateRequest) (*shortrpc.CreateResponse, error) {

	res := shortrpc.CreateResponse{
		Status: "OK",
	}
	log.Info(in.GetOriginalUrl())
	hashcode, err := storage.ParserDataURL(in.GetOriginalUrl())

	if err != nil {
		log.Error(err)
		return &res, status.Error(codes.FailedPrecondition, err.Error())
	}

	log.Info(hashcode.ShortPath)
	hashcode.UserID = in.GetUserId()
	gl, err := hook.Storage.Put(hashcode.ShortPath, hashcode)

	if err != nil {
		log.Error(err)
		res.Status = err.Error()
		return &res, status.Error(codes.FailedPrecondition, err.Error())

	}
	res.ResponseUrl = hashcode.ShortPath
	if len(gl) > 0 {
		res.ResponseUrl = gl
	}
	return &res, nil
}
