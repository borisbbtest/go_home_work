package handlersgrpc

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/proto/shortrpc"
	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

func (hook *WrapperHandlerRPC) GetStats(ctx context.Context, in *shortrpc.GetStatsRequest) (*shortrpc.GetStatsResponse, error) {
	res := shortrpc.GetStatsResponse{
		Status: "Ok",
		Code:   1,
	}

	_, err := tools.TrustedSubnetIP(in.GetIpAddress(), *hook.ServerConf.Subnet)
	if err != nil {
		res.Status = err.Error()
		return &res, status.Error(codes.PermissionDenied, "Please, You will call the administrator for access in endpoint")
	}

	value, err := hook.Storage.GetStats()

	if err != nil {
		res.Status = err.Error()
		return &res, status.Error(codes.FailedPrecondition, "Opps OMG")
	}
	res.Urls = value.URLs
	res.Users = value.Users
	return &res, status.Error(codes.OK, "It is good response")
}
