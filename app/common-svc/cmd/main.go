package main

import (
	"fortune-bd/app/common-svc/internal/service"
	"fortune-bd/app/common-svc/server"
	"fortune-bd/libs/env"
	"github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"os"
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
	grpcServers := server.NewGRPCServers(service.NewCommonService())
	httpServer := server.NewHTTPServer()
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(env.EXCHANGE_SRV_NAME),
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