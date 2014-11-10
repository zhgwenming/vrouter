package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
)

type Command struct {
	etcdConfig *registry.ClientConfig
}

func NewCommand(cfg *registry.ClientConfig) *Command {
	c := &Command{cfg}
	return c
}

func (c *Command) CellInit(parent *cobra.Command) {
	cell := new(Cell)
	cell.cmd = c
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

func (c *Command) Service(parent *cobra.Command) {
	srv := new(Service)
	srv.cmd = c
	// new subcommand
	cmd := &cobra.Command{
		Use:   "service [add|delete]",
		Short: "service management",
		Long:  "",
		Run:   srv.Run,
	}

	cmd.Flags().StringVarP(&srv.Name, "name", "n", srv.Name, "service name")
	cmd.Flags().StringVarP(&srv.Addr, "listen", "l", srv.Addr, "service listen address")
	cmd.Flags().StringVarP(&srv.Port, "port", "p", srv.Port, "service port")

	parent.AddCommand(cmd)
}
