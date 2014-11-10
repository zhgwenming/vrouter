package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
)

type cli struct {
	etcdConfig *registry.ClientConfig
}

func NewCli(cfg *registry.ClientConfig) *cli {
	c := &cli{cfg}
	return c
}

func (c *cli) CellInit(parent *cobra.Command) {
	cell := new(Cell)
	cell.cli = c
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

func (c *cli) Service(parent *cobra.Command) {
}
