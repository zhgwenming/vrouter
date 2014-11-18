package daemon

import (
	daemonctl "github.com/zhgwenming/gbalancer/daemon"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/netinfo"
	"github.com/zhgwenming/vrouter/registry"
	"log"
	"net"
	"os"
	"runtime"
)

// command for configuration which exists before and after *cobra.Config.Execute()
type Config struct {
	// command switches
	daemonMode  bool
	gatewayMode bool
	foreground  bool

	pidFile string
	// host relate information
	Hostip   string
	Hostname string

	// bridge information
	BridgeName string

	etcdConfig *registry.ClientConfig

	// vrouter daemon
	Daemon *Daemon
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) InitCmd(client *registry.ClientConfig) *cobra.Command {

	cfg.etcdConfig = client

	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run:  cfg.Run,
	}

	var ipnet *net.IPNet
	var hostIp string

	ipnetlist := netinfo.ListIPNet(true)
	if len(ipnetlist) > 0 {
		ipnet = ipnetlist[0]
		hostIp = ipnet.String()
	} else {
		hostIp = ""
	}

	// vrouter flags
	flags := routerCmd.Flags()

	cfg.Hostname, _ = os.Hostname()

	flags.BoolVarP(&cfg.daemonMode, "daemon", "d", false, "whether to run as daemon mode")
	flags.BoolVarP(&cfg.foreground, "foreground", "f", false, "whether to run as a foreground process")
	flags.StringVarP(&cfg.pidFile, "pidfile", "p", "", "pidfile to write to")

	flags.BoolVarP(&cfg.gatewayMode, "gateway", "g", false, "to run as dedicated gateway, will not allocate subnet on this machine")

	// need to convert to IPNet form
	flags.StringVarP(&cfg.Hostip, "hostip", "i", hostIp, "use specified ip/mask instead auto detected ip address")

	// vrouter information
	flags.StringVarP(&cfg.Hostname, "hostname", "n", cfg.Hostname, "hostname to use in daemon mode")
	flags.StringVarP(&cfg.BridgeName, "bridge", "b", "docker0", "bridge name to setup")

	return routerCmd
}

// Serve is the actual entry point of the handler
func (cfg *Config) Serve() {
	client := registry.NewClient(cfg.etcdConfig)
	vrouter := NewDaemon(cfg, client)
	cfg.Daemon = vrouter

	// start keepalive first
	err := vrouter.KeepAlive()
	if err != nil {
		log.Fatalf("error to keepalive: %s", err)
	}

	// bind and get a bridge IPNet with our iface ip
	// create the routing table entry in registry
	bridgeIPNet, err := vrouter.BindBridgeIPNet(cfg.Hostip)
	if err != nil {
		log.Fatal("Failed to bind router interface: ", err)
	} else {
		log.Printf("Requested bridge ip - %v\n", bridgeIPNet)
	}

	// create bridge if we're running under linux
	// to debug on Mac OS X
	if runtime.GOOS == "linux" {
		err = vrouter.CreateBridge(bridgeIPNet.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {

		// monitor the routing table change
		err = vrouter.ManageRoute()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (cfg *Config) Run(c *cobra.Command, args []string) {
	if cfg.daemonMode {
		//daemon := cmd.Daemon
		// -peer-addr 127.0.0.1:7001 -addr 127.0.0.1:4001 -data-dir machines/machine1 -name machine1
		//go registry.StartEtcd("-peer-addr", "127.0.0.1:7001", "-addr", "127.0.0.1:4001", "-data-dir", "machines/"+daemon.Hostname, "-name", daemon.Hostname)

		daemonctl.HandleFunc(cfg.Serve)

		// start the engine
		if err := daemonctl.Start(cfg.pidFile, cfg.foreground); err != nil {
			log.Fatal(err)
		}

	} else {
		c.Help()
	}
}
