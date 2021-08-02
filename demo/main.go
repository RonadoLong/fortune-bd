package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	rpcpb "wq-fotune-backend/demo/proto"
	"wq-fotune-backend/pkg/etcd"
)

const (
	port = "127.0.0.1:9992"
	srv_name = "demo-srv"
	reg = "http://127.0.0.1:2379"
)

var (
	errRPCCancel = errors.New("rpc cancel")
)

type DemoHandler struct {

}

func (d DemoHandler) Hello(ctx context.Context, req *rpcpb.Req) (*rpcpb.Response, error) {
	log.Println(req.Id)
	select {
	case<-ctx.Done():
		return nil, errRPCCancel
	default:
		return &rpcpb.Response{
			Resp: fmt.Sprintf("hello %s", req.Id),
		}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	err = etcd.RegisterServer(etcd.NewClient(reg), srv_name, port, nil)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	rpcpb.RegisterDemoServer(s, &DemoHandler{})
	go waitStop(func() {
		s.GracefulStop()
	})
	log.Printf("starting hello internal at %s", port)
	if err = s.Serve(lis); err != nil{
		log.Fatal(err)
	}
}

func waitStop(exit func()) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sc
	log.Printf("exit: signal=<%d>. \n", sig)
	exit()
	switch sig {
	case syscall.SIGTERM:
		log.Println("exit: bye :-).")
		os.Exit(0)
	default:
		log.Println("exit: bye :-(.")
		os.Exit(1)
	}
}
