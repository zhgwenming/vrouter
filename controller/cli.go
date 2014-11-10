package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
	"reflect"
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

// InitCmd creates etcd client and register cobra subcommand
func InitCmd(parent *cobra.Command, client *registry.ClientConfig) {
	etcdConfig = client
	c := new(cli)

	v := reflect.ValueOf(c)
	for i := 0; i < v.NumMethod(); i++ {
		subcmd := v.Method(i).Interface().(func(*cobra.Command))
		subcmd(parent)
	}
}
