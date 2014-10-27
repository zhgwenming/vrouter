package daemon

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/netinfo"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

type Command struct {
	etcdServers *string
	hostip      string

	// command switches
	daemonMode  bool
	gatewayMode bool

	// vrouter daemon
	daemon *Daemon
}

func NewCommand() *Command {
	return &Command{}
}

func (cmd *Command) InitCmd(servers *string) *cobra.Command {

	vrouter := NewDaemon()
	cmd.daemon = vrouter

	cmd.etcdServers = servers

	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run:  cmd.Run,
	}

	var ipnet *net.IPNet
	ipnetlist := netinfo.ListIPNet(true)
	if len(ipnetlist) > 0 {
		ipnet = ipnetlist[0]
	}

	// vrouter flags
	cmdflags := routerCmd.Flags()

	vrouter.Hostname, _ = os.Hostname()

	cmdflags.BoolVarP(&cmd.daemonMode, "daemon", "d", false, "whether to run as daemon mode")
	cmdflags.BoolVarP(&cmd.gatewayMode, "gateway", "g", false, "to run as dedicated gateway, will not allocate subnet on this machine")

	// need to convert to IPNet form
	cmdflags.StringVarP(&cmd.hostip, "hostip", "i", ipnet.String(), "use specified ip/mask instead auto detected ip address")

	// vrouter information
	cmdflags.StringVarP(&vrouter.Hostname, "hostname", "n", vrouter.Hostname, "hostname to use in daemon mode")
	cmdflags.StringVarP(&vrouter.bridgeName, "bridge", "b", "docker0", "bridge name to setup")

	return routerCmd
}

func (cmd *Command) Run(c *cobra.Command, args []string) {
	if cmd.daemonMode {
		servers := strings.Split(*cmd.etcdServers, ",")
		vrouter := cmd.daemon

		// start keepalive first
		vrouter.etcdClient = etcd.NewClient(servers)
		err := vrouter.KeepAlive()
		if err != nil {
			log.Fatalf("error to keepalive: %s, other instance running?", err)
		}

		// bind and get a bridge IPNet with our iface ip
		// create the routing table entry in registry
		bridgeIPNet, err := vrouter.BindBridgeIPNet(cmd.hostip)
		if err != nil {
			log.Fatal("Failed to bind router interface: ", err)
		} else {
			log.Printf("daemon: get ipnet %v\n", bridgeIPNet)
		}

		// create bridge if we're running under linux
		// to debug on Mac OS X
		if runtime.GOOS == "linux" {
			err = vrouter.CreateBridge(bridgeIPNet.String())
			if err != nil {
				log.Fatal(err)
			}
		}

		// monitor the routing table change
		err = vrouter.ManageRoute()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		c.Help()
	}
}
