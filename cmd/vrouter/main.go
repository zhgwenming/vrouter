package main

import (
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/pkg/registry"
	"log"
)

var (
	daemon     bool
	etcdServer string
)

func main() {
	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	routerCmd.Flags().BoolVarP(&daemon, "daemon", "d", true, "whether to run as daemon mode")
	routerCmd.PersistentFlags().StringVarP(&etcdServer, "etcd_server", "e", "127.0.0.1:4001", "etcd registry addr")

	registry.Init(routerCmd, &etcdServer)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
