package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	daemon bool
	subnet string
)

const (
	DEFAULT_SUBNET = "10.0.0.0/16"
)

func registryInit(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		subnet = args[0]
	} else {
		subnet = DEFAULT_SUBNET
	}

	fmt.Printf("vrouter init %s\n", subnet)
}

func main() {
	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	routerCmd.Flags().BoolVarP(&daemon, "daemon", "d", true, "whether to run as daemon mode")

	initCmd := &cobra.Command{
		Use:   "init [subnet]",
		Short: "init the registry",
		Long:  "init the registry with speciffic ip network information",
		Run:   registryInit,
	}

	routerCmd.AddCommand(initCmd)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
