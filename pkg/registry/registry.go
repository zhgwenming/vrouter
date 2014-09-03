package registry

import (
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"log"
	"net"
	"strings"
)

var (
	machines       string
	hostNames      []string
	globalSubnet   string
	etcdServer     []string
	registryClient *Registry
)

const (
	DEFAULT_SUBNET  = "10.0.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

type Registry struct {
	etcdClient *etcd.Client
}

func NewRegistry(etcdClient *etcd.Client) *Registry {
	return &Registry{etcdClient: etcdClient}
}

// create etcd client
// register cobra subcommand
func Init(parent *cobra.Command, etcdServerStr string) {

	etcdServer = strings.Split(etcdServerStr, ",")

	etcdClient := etcd.NewClient(etcdServer)
	registryClient = NewRegistry(etcdClient)

	// register new subcommand
	initCmd := &cobra.Command{
		Use:   "init <machine1,machine2,..>",
		Short: "init the machine registry",
		Long:  "init the machine registry with specific ip network information",
		Run:   registryInit,
	}

	initCmd.Flags().StringVarP(&globalSubnet, "ipnet", "n", DEFAULT_SUBNET, "cidr ip subnet information")

	parent.AddCommand(initCmd)
}

func registryRoutePrefix() string {
	return REGISTRY_PREFIX + "/" + "route"
}

func registryInit(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		machines = args[0]
		// just make slice if machines is not empty
		if len(machines) > 0 {
			hostNames = strings.Split(machines, ",")
		} else {
			cmd.Help()
			log.Fatal("Empty machine list specified")
		}
	} else {
		cmd.Help()
		log.Fatal("No machine list specified")
	}

	_, ipnet, err := net.ParseCIDR(globalSubnet)

	if err != nil {
		log.Fatal(err)
	}

	nets := GetAllSubnet(ipnet, 8)

	fmt.Printf("vrouter init %s, %v, etcd: %s\n", globalSubnet, ipnet, etcdServer)
	//fmt.Printf("%v\n", nets)

	routePrefix := registryRoutePrefix()

	//fmt.Printf("hostnames %d, %v\n", len(hostNames), hostNames)
	for i, node := range hostNames {
		key := routePrefix + "/" + node + "/" + "ipnet"
		log.Printf("initialize config for host %s\n", node)
		if _, err := registryClient.etcdClient.Create(key, nets[i].String(), 0); err != nil {
			log.Printf("Error to create node: %s", err)
		}

	}
}
