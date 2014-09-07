package daemon

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/netinfo"
	"log"
	"net"
	"os"
	"strings"
)

type Command struct {
	etcdServers *string
	daemonMode  bool
	gatewayMode bool
	hostname    string
	hostip      string
	bridgeName  string
}

func NewCommand() *Command {
	return &Command{}
}

func (cmd *Command) InitCmd(servers *string) *cobra.Command {

	cmd.etcdServers = servers

	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run:  Run,
	}

	var ipnet *net.IPNet
	ipnetlist := netinfo.ListIPNet(true)
	if len(ipnetlist) > 0 {
		ipnet = ipnetlist[0]
	}

	// vrouter flags
	cmdflags := routerCmd.Flags()

	hostname, _ = os.Hostname()

	cmdflags.BoolVarP(&cmd.daemonMode, "daemon", "d", false, "whether to run as daemon mode")
	cmdflags.BoolVarP(&cmd.gatewayMode, "gateway", "g", false, "to run as dedicated gateway, will not allocate subnet on this machine")
	cmdflags.StringVarP(&cmd.hostname, "hostname", "n", hostname, "hostname to use in daemon mode")
	cmdflags.StringVarP(&cmd.hostip, "hostip", "i", ipnet.String(), "use specified ip/mask instead auto detected ip address")
	cmdflags.StringVarP(&cmd.bridgeName, "bridge", "b", "docker0", "bridge name to setup")

	return routerCmd
}

func (cmd *Command) Run(c *cobra.Command, args []string) {
	if daemonMode {
		servers := strings.Split(*cmd.etcdServers, ",")
		etcdClient := etcd.NewClient(servers)
		vrouter = NewDaemon(etcdClient)

		vrouter.KeepAlive(hostname)
		dockerNet, err := BindDockerNet(hostname, hostip)
		if err != nil {
			log.Fatal("Failed to bind router interface: ", err)
		} else {
			log.Printf("daemon: get ipnet %v\n", dockerNet)
		}

		err = vrouter.createBridgeIface(bridge, dockerNet.String())
		if err != nil {
			log.Fatal(err)
		}

		err = vrouter.ManageRoute()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		c.Help()
	}
}
