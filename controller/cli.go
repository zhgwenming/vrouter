package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
)

type cli struct{}

func (c *cli) Init(parent *cobra.Command) {
	// new subcommand
	initCmd := &cobra.Command{
		Use:   "init <machine1,machine2,..>",
		Short: "init the machine registry",
		Long:  "init the machine registry with specific ip network information",
		Run:   registryInit,
	}

	initCmd.Flags().StringVarP(&cellSubnet, "cellnet", "c", registry.DEFAULT_SUBNET, "cell cidr subnet ip address")
	initCmd.Flags().StringVarP(&overlaySubnet, "overlay", "o", registry.DEFAULT_SUBNET, "the whole overlay subnet ip address")

	parent.AddCommand(initCmd)
}

func (c *cli) Service(parent *cobra.Command) {
}
