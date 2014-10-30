package registry

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type ClientConfig struct {
	Servers string

	// tls authentication related
	CaFile   string
	CertFile string
	KeyFile  string
}

func NewClient(cfg *ClientConfig) *etcd.Client {
	var client *etcd.Client
	var err error

	//if cmd.CaFile != "" && cmd.CertFile != "" && cmd.KeyFile != "" {
	servers := strings.Split(cfg.Servers, ",")
	log.Printf("%v", servers)
	if strings.HasPrefix(servers[0], "https://") {
		client, err = etcd.NewTLSClient(servers, cfg.CertFile, cfg.KeyFile, cfg.CaFile)
		if err != nil {
			log.Fatalf("error to create tls client: %s", err)
		}

		log.Printf("established tls connection.")
	} else {
		client = etcd.NewClient(servers)
		log.Printf("established plain text connection.")
	}
	log.Printf("Created new etcd client")
	return client
}
