package main

import (
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/pkg/registry"
	"log"
)

var (
	daemon bool
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

	registry.Init(routerCmd)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
