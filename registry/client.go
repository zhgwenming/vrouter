package registry

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type ClientConfig struct {
	Servers string
	Verbose bool

	// tls authentication related
	CaFile   string
	CertFile string
	KeyFile  string
}

func NewClient(cfg *ClientConfig) *etcd.Client {
	var client *etcd.Client
	var err error

	servers := strings.Split(cfg.Servers, ",")
	if cfg.Verbose {
		log.Printf("etcd client with: %v", servers)
	}

	// cert and key file are needed for tls authentication
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		client, err = etcd.NewTLSClient(servers, cfg.CertFile, cfg.KeyFile, cfg.CaFile)
		if err != nil {
			log.Fatalf("error to create tls client: %s", err)
		}

		if cfg.Verbose {
			log.Printf("established tls connection.")
		}
	} else {
		client = etcd.NewClient(servers)
		if cfg.Verbose {
			log.Printf("established plain text connection.")
		}
	}
	if cfg.Verbose {
		log.Printf("Created new etcd client")
	}
	return client
}
