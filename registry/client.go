package registry

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"log"
)

type ClientConfig struct {
	EtcdServers string

	// tls authentication related
	CaFile   string
	CertFile string
	KeyFile  string
}

func NewClient(cfg *ClientConfig) *etcd.Client {

	log.Printf("Creating new etcd client")
	return nil
}
