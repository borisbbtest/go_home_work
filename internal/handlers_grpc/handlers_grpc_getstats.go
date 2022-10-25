package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/tools"
)

func (hook *WrapperHandlerRPC) GetStats(ctx context.Context, in *shortrpc.GetStatsRequest) (*shortrpc.GetStatsResponse, error) {
	res := shortrpc.GetStatsResponse{
		Status: "ok",
	}

	_, err := tools.TrustedSubnetIP(in.GetIpAddress(), hook.ServerConf.TrustedSubnet)
	if err != nil {
		res.Status = err.Error()
		return &res, nil
	}

	value, err := hook.Storage.GetStats()

	if err != nil {
		res.Status = err.Error()
	}
	res.Urls = value.URLs
	res.Users = value.Users
	return &res, nil
}
