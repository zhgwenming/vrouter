package registry

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var (
	machines    string
	hostNames   []string
	subnet      string
	etcdServers []string
	client      etcdRegistry
)

const (
	DEFAULT_SUBNET = "10.0.0.0/16"
)

func Init(parent *cobra.Command, etcd string) {

	etcdServers = strings.Split(etcd, ",")

	initCmd := &cobra.Command{
		Use:   "init [subnet]",
		Short: "init the registry",
		Long:  "init the registry with speciffic ip network information",
		Run:   registryInit,
	}

	initCmd.Flags().StringVarP(&machines, "machines", "m", "", "List of machines")

	parent.AddCommand(initCmd)
}

type etcdRegistry struct {
	etcdClient *etcd.Client
}

func registryInit(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		subnet = args[0]
	} else {
		subnet = DEFAULT_SUBNET
	}

	fmt.Printf("vrouter init %s, etcd: %s\n", subnet, etcdServers)
	etcd := etcd.NewClient(etcdServers)
	client = etcdRegistry{etcdClient: etcd}
}

func (r *etcdRegistry) Set(key, value string, ttl uint64) error {
	client := r.etcdClient

	if resp, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	} else {
		//log.Printf("No instance exist on this node, starting")
		go func() {
			sleeptime := time.Duration(ttl / 3)
			for {
				index := resp.EtcdIndex
				time.Sleep(sleeptime * time.Second)
				resp, err = client.CompareAndSwap(key, value, ttl, value, index)
				if err != nil {
					log.Fatal("Unexpected lost our node lock", err)
				}
			}
		}()
		return nil
	}
}
