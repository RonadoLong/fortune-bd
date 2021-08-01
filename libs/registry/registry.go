package registry

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

var ir registry.Registry

// GetIns etcd注册中心
func GetIns(etcdAddr string) registry.Registry {
	if ir == nil {
		ir = etcd.NewRegistry(registry.Addrs(etcdAddr))
	}
	return ir
}
