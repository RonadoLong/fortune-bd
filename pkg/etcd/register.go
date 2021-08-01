package etcd

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/naming"
	"log"
	"time"
)

var (
	timeout = time.Second * 10
	ttl     = int64(15)
	Prefix  = "ifortune-services"
)

func NewClient(target ...string) *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: target,
	})
	if err != nil {
		log.Println(fmt.Errorf("grpclb: create etcd3 client failed: %v", err))
		return nil
	}
	return client
}

func RegisterServer(client *clientv3.Client, serviceName, addr string, metadata interface{}) error {
	metaData := naming.Update{
		Op:       naming.Add,
		Addr:     adjustAddr(addr),
		Metadata: metadata,
	}
	publisher, err := newEtcdPublisher(client, Prefix, ttl, timeout)
	if err != nil {
		return err
	}
	err = publisher.Publish(serviceName, metaData)
	if err != nil {
		log.Fatalf("rpc: publish service <%s> failed, error:\n%+v", serviceName, err)
		return err
	}
	log.Printf("rpc: service <%s> already published \n", serviceName)
	return nil
}

func adjustAddr(addr string) string {
	if addr[0] == ':' {
		ips, err := intranetIP()
		if err != nil {
			log.Fatalf("get intranet ip failed, error:\n%+v", err)
		}
		return fmt.Sprintf("%s%s", ips[0], addr)
	}
	return addr
}

// UnRegister delete registered service from etcd
//func UnRegister() error {
//	stopSignal <- true
//	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
//	var err error
//	if _, err := client.Delete(context.Background(), serviceKey); err != nil {
//		log.Printf("grpclb: deregister '%s' failed: %s", serviceKey, err.Error())
//	} else {
//		log.Printf("grpclb: deregister '%s' ok.", serviceKey)
//	}
//	return err
//}
