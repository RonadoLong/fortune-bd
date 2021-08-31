package main

import (
	"context"
	"fortune-bd/app/usercenter-svc/internal/service"
	"fortune-bd/app/usercenter-svc/server"
	"fortune-bd/libs/env"
	"fortune-bd/libs/logger"
	"github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	id, _ = os.Hostname()
)


func main() {
	log.Println("id:", id)
	client, err := etcd.New(etcd.Config{
		Endpoints: []string{env.EtcdAddr},
	})
	if err != nil {
		log.Fatal(err)
	}
	r := registry.New(client)
	grpcServers := server.NewGRPCServers(service.NewUserService())
	httpServer := server.NewHTTPServer()
	defer func() {
		grpcServers.GracefulStop()
		httpServer.Shutdown(context.Background())
	}()
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(env.USER_SRV_NAME),
		kratos.Version("1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			httpServer,
			grpcServers,
		),
		kratos.Registrar(r),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	<-c
	logger.Info("服务已关闭")
}