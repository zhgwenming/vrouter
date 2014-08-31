package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	client bool
)

func main() {
	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	routerCmd.Flags().BoolVarP(&client, "client", "c", true, "whether to run as client mode")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init the registry",
		Run: func(c *cobra.Command, args []string) {
			fmt.Printf("vrouter init\n")
		},
	}

	routerCmd.AddCommand(initCmd)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
