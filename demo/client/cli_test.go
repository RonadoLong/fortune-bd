package client

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"testing"
	"time"
)

func TestNewCli(t *testing.T) {
	NewCli()
}

func TestNewCli2(t *testing.T) {
	cli, err := clientv3.NewFromURL("127.0.0.1:2379")
	if err != nil{
		log.Println(err)
		return
	}
	response, err := cli.Put(context.Background(), "long", "1111111111")
	log.Println(response, err)

	for t := range time.NewTicker(time.Second).C {
		resp, err := cli.Get(context.Background(), srv_name)
		log.Println(err)
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Kvs[0].Value)
		}
	}

}