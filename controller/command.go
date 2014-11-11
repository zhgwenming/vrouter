package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
)

type Config struct {
	etcdConfig *registry.ClientConfig
}

func NewConfig(cfg *registry.ClientConfig) *Config {
	c := &Config{cfg}
	return c
}

func (c *Config) CellInit(parent *cobra.Command) {
	cell := new(NodeManager)
	cell.config = c
	// new subcommand
	cmd := &cobra.Command{
		Use:   "cell-init <machine1,machine2,..>",
		Short: "cell-init the machine registry",
		Long:  "init the machine registry with specific ip network information",
		Run:   cell.registryInit,
	}

	cmd.Flags().StringVarP(&cell.cellSubnet, "cellnet", "c", registry.DEFAULT_SUBNET, "cell cidr subnet ip address")
	cmd.Flags().StringVarP(&cell.overlaySubnet, "overlay", "o", registry.DEFAULT_SUBNET, "the whole overlay subnet ip address")

	parent.AddCommand(cmd)
}

func (c *Config) Service(parent *cobra.Command) {
	manager := new(ServiceManager)
	manager.config = c
	// new subcommand
	cmd := &cobra.Command{
		Use:   "service [list|add|delete]",
		Short: "service management",
		Long:  "manage the distributed network services",
		Run:   manager.Run,
	}

	cmd.Flags().StringVarP(&manager.srvConfig.Name, "name", "n", "", "service name")
	cmd.Flags().StringVarP(&manager.srvConfig.Addr, "listen", "l", "", "service listen address")
	cmd.Flags().StringVarP(&manager.srvConfig.Port, "port", "p", "", "service port")
	cmd.Flags().StringVarP(&manager.srvConfig.BackEnds, "backend", "b", "", "'ip1:port1 ip2:port2' form of backends")

	parent.AddCommand(cmd)
}
