package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	 etcdnaming "go.etcd.io/etcd/clientv3/naming"
	"google.golang.org/grpc/naming"
	"log"
	"time"
)

// Publisher a service publish
type Publisher interface {
	Publish(service string, meta naming.Update) error
}

type etcdPublisher struct {
	prefix   string
	ttl      int64
	timeout  time.Duration
	client   *clientv3.Client
	resolver *etcdnaming.GRPCResolver
}

func newEtcdPublisher(client *clientv3.Client, prefix string, ttl int64, timeout time.Duration) (Publisher, error) {
	return &etcdPublisher{
		prefix:  prefix,
		ttl:     ttl,
		timeout: timeout,
		client:  client,
		resolver: &etcdnaming.GRPCResolver{
			Client: client,
		},
	}, nil
}

func (p *etcdPublisher) Publish(service string, meta naming.Update) error {
	lessor := clientv3.NewLease(p.client)
	defer lessor.Close()
	ctx, cancel := context.WithTimeout(p.client.Ctx(), p.timeout)
	leaseResp, err := lessor.Grant(ctx, p.ttl)
	cancel()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = p.client.KeepAlive(p.client.Ctx(), leaseResp.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, cancel = context.WithTimeout(p.client.Ctx(), p.timeout)
	defer cancel()
	return p.resolver.Update(ctx, fmt.Sprintf("%s/%s", p.prefix, service), meta, clientv3.WithLease(clientv3.LeaseID(leaseResp.ID)))
}

