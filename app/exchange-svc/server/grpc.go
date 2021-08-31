package server

import (
	v1 "fortune-bd/api/exchange/v1"
	"fortune-bd/app/exchange-svc/internal/service"
	"github.com/go-kratos/kratos/middleware/recovery/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
)


// NewGRPCServers new a gRPC server.
func NewGRPCServers (service *service.ExOrderService) *grpc.Server {
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
	v1.RegisterExOrderServer(srv, service)
	return srv
}