package etcd

import (
	"errors"
	"fmt"
	etcd3 "go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/naming"
	"strings"
)

// resolver is the implementaion of grpc.naming.Resolver
type resolver struct {
	serviceName string // internal name to resolve
}

// NewResolver return resolver with internal name
func NewResolver(serviceName string) *resolver {
	return &resolver{serviceName: serviceName}
}

// Resolve to resolve the internal from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == "" {
		return nil, errors.New("grpclb: no internal name provided")
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed: %s", err.Error())
	}
	// Return watcher
	return &watcher{re: re, client: *client}, nil
}
