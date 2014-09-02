package registry

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"log"
	"net"
	"strings"
	"time"
)

var (
	machines     string
	hostNames    []string
	globalSubnet string
	etcdServers  []string
	client       etcdRegistry
)

const (
	DEFAULT_SUBNET  = "10.0.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

func Init(parent *cobra.Command, etcd string) {

	etcdServers = strings.Split(etcd, ",")

	initCmd := &cobra.Command{
		Use:   "init <machine1,machine2,..>",
		Short: "init the machine registry",
		Long:  "init the machine registry with specific ip network information",
		Run:   registryInit,
	}

	initCmd.Flags().StringVarP(&globalSubnet, "ipnet", "n", DEFAULT_SUBNET, "cidr ip subnet information")

	parent.AddCommand(initCmd)
}

type etcdRegistry struct {
	etcdClient *etcd.Client
}

func GetAllSubnet(ipnet *net.IPNet, hostbits int) []net.IPNet {
	ones, bits := ipnet.Mask.Size()
	zeros := bits - ones

	// network bits
	netBits := zeros - hostbits
	if netBits < 0 {
		return []net.IPNet{}
	}

	ip4 := ipnet.IP.To4()

	numberSubnet := 1 << uint(netBits)
	subnet := make([]net.IPNet, 0, numberSubnet)

	for i := uint32(0); i < uint32(numberSubnet); i++ {
		ipbuf := make([]byte, 4)
		number := i << uint(hostbits)
		binary.BigEndian.PutUint32(ipbuf, number)

		ip := (((uint32(ipbuf[0]) | uint32(ip4[0])) << 24) |
			((uint32(ipbuf[1]) | uint32(ip4[1])) << 16) |
			((uint32(ipbuf[2]) | uint32(ip4[2])) << 8) |
			uint32(ipbuf[3]) | uint32(ip4[3]))
		binary.BigEndian.PutUint32(ipbuf, ip)

		ipmask := net.CIDRMask(bits-hostbits, bits)

		subipnet := net.IPNet{ipbuf, ipmask}
		subnet = append(subnet, subipnet)
	}

	return subnet

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

	fmt.Printf("vrouter init %s, %v, etcd: %s\n", globalSubnet, ipnet, etcdServers)
	//fmt.Printf("%v\n", nets)
	etcd := etcd.NewClient(etcdServers)
	client = etcdRegistry{etcdClient: etcd}

	keyPrefix := REGISTRY_PREFIX + "/" + "route"
	//fmt.Printf("hostnames %d, %v\n", len(hostNames), hostNames)
	for i, node := range hostNames {
		key := keyPrefix + "/" + node + "/" + "ipnet"
		log.Printf("initialize config for host %s\n", node)
		if value, err := json.Marshal(nets[i]); err != nil {
			log.Fatal(err)
		} else {
			if _, err := etcd.Create(key, string(value), 0); err != nil {
				log.Printf("Error to create node: %s", err)
			}
		}
	}
}

func (r *etcdRegistry) KeepAlive(key, value string, ttl uint64) error {
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
