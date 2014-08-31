package main

import (
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
	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
