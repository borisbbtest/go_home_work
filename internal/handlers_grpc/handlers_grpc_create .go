package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/storage"
)

func (hook *WrapperHandlerRPC) Create(ctx context.Context, in *shortrpc.CreateRequest) (*shortrpc.CreateResponse, error) {

	res := shortrpc.CreateResponse{
		Status: "OK",
	}
	log.Info(in.GetOriginalUrl())
	hashcode, err := storage.ParserDataURL(in.GetOriginalUrl())

	if err != nil {
		log.Info(err)
	}

	log.Info(hashcode.ShortPath)
	hashcode.UserID = in.GetUserId()
	gl, err := hook.Storage.Put(hashcode.ShortPath, hashcode)

	if err != nil {
		log.Error(err)
		res.Status = err.Error()
	}
	res.ResponseUrl = hashcode.ShortPath
	if len(gl) > 0 {
		res.ResponseUrl = gl
	}
	return &res, nil
}
