package controller

import (
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/netinfo"
	"github.com/zhgwenming/vrouter/registry"
	"log"
	"net"
	"strings"
)

type Cell struct {
	cmd           *Command
	machines      string
	hostNames     []string
	cellSubnet    string
	overlaySubnet string
	etcdClient    *etcd.Client
}

func (n *Cell) registryInit(cmd *cobra.Command, args []string) {

	n.etcdClient = registry.NewClient(n.cmd.etcdConfig)

	if len(args) > 0 {
		n.machines = args[0]
		// just make slice if machines is not empty
		if len(n.machines) > 0 {
			n.hostNames = strings.Split(n.machines, ",")
		} else {
			cmd.Help()
			log.Fatal("Empty machine list specified")
		}
	} else {
		cmd.Help()
		log.Fatal("No machine list specified")
	}

	_, ipnet, err := net.ParseCIDR(n.cellSubnet)

	if err != nil {
		log.Fatal(err)
	}

	nets := netinfo.GetAllSubnet(ipnet, 8)

	fmt.Printf("vrouter init %s, %v\n", n.cellSubnet, ipnet)
	//fmt.Printf("%v\n", nets)

	//fmt.Printf("hostnames %d, %v\n", len(hostNames), hostNames)
	for i, node := range n.hostNames {
		key := registry.BridgeInfoPath(node)
		log.Printf("initialize config for host %s\n", node)
		if _, err := n.etcdClient.Create(key, nets[i].String(), 0); err != nil {
			log.Printf("Error to create node: %s", err)
		}

	}

	// create the overlay network ip info
	key := registry.RouterOverlayPath()
	log.Printf("initialize overlay network ip information with %s\n", n.overlaySubnet)
	if _, err := n.etcdClient.Create(key, n.overlaySubnet, 0); err != nil {
		log.Printf("Error to create node: %s", err)
	}
}
