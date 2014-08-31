package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

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
