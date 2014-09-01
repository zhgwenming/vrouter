package registry

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	machines   string
	hostNames  []string
	subnet     string
	etcdServer *string
)

const (
	DEFAULT_SUBNET = "10.0.0.0/16"
)

func Init(parent *cobra.Command, etcd *string) {

	etcdServer = etcd

	initCmd := &cobra.Command{
		Use:   "init [subnet]",
		Short: "init the registry",
		Long:  "init the registry with speciffic ip network information",
		Run:   registryInit,
	}

	initCmd.Flags().StringVarP(&machines, "machines", "m", "", "List of machines")

	parent.AddCommand(initCmd)
}

func registryInit(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		subnet = args[0]
	} else {
		subnet = DEFAULT_SUBNET
	}

	fmt.Printf("vrouter init %s, etcd: %s\n", subnet, *etcdServer)
}
