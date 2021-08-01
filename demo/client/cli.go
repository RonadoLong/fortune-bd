package client

import (
	"context"
	"fmt"
	etcdnaming "go.etcd.io/etcd/clientv3/naming"
	"google.golang.org/grpc"
	"log"
	"time"
	rpcpb "wq-fotune-backend/demo/proto"
	"wq-fotune-backend/pkg/etcd"
)

const (
	port     = 9991
	srv_name = "demo-srv"
	reg      = "http://127.0.0.1:2379"
)

// NewCli
func NewCli() {
	//r := etcd.NewResolver(fmt.Sprintf("%s/%s", etcd.Prefix, srv_name))
	//watcher, err2 := r.Resolve("127.0.0.1:2379")
	//if err2 != nil {
	//	log.Println(err2)
	//	return
	//}
	c := etcd.NewClient(reg)
	r := &etcdnaming.GRPCResolver{
		Client: c,
	}

	var grpcOptions []grpc.DialOption
	grpcOptions = append(grpcOptions, grpc.WithInsecure())
	grpcOptions = append(grpcOptions, grpc.WithTimeout(time.Second*10))
	grpcOptions = append(grpcOptions, grpc.WithBlock())
	grpcOptions = append(grpcOptions, grpc.WithBalancer(grpc.RoundRobin(r)))

	conn, err := grpc.Dial(fmt.Sprintf("%s/%s", etcd.Prefix, srv_name), grpcOptions...)
	if err != nil {
		log.Println(err)
		return
	}
	client := rpcpb.NewDemoClient(conn)
	for t := range time.NewTicker(time.Second).C {
		resp, err := client.Hello(context.Background(), &rpcpb.Req{Id: "22222"})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%v: Reply is %s\n", t, resp.Resp)
	}
}
