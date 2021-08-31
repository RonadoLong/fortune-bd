package server

import (
	v1 "fortune-bd/api/usercenter/v1"
	"fortune-bd/app/usercenter-svc/internal/service"
	"github.com/go-kratos/kratos/middleware/recovery/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
)


// NewGRPCServers new a gRPC server.
func NewGRPCServers (service *service.UserService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			//tracing.Server(
			//	tracing.WithTracerProvider(tp)),
			//logging.Server(logger),
		),
	}
	opts = append(opts, grpc.Timeout(time.Second* 5), grpc.Address(":9000"))
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServer(srv, service)
	return srv
}