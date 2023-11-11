package grpc

import (
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	userservice "harmoni/app/harmoni/internal/service/user"
	"harmoni/internal/conf"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGrpcServer(
	conf *conf.Server,
	logger log.Logger,
	userservice *userservice.UserGRPCService,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	}
	if conf.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(conf.Grpc.Addr))
	}
	if conf.Grpc.Timeout != 0 {
		opts = append(opts, grpc.Timeout(conf.Grpc.Timeout))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServer(srv, userservice)
	return srv
}
