package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
	"reflect"
)

// AddSubCommands creates etcd client and register cobra subcommand
func AddSubCommands(parent *cobra.Command, client *registry.ClientConfig) {
	c := NewCommand(client)

	v := reflect.ValueOf(c)
	for i := 0; i < v.NumMethod(); i++ {
		subcmd := v.Method(i).Interface().(func(*cobra.Command))
		subcmd(parent)
	}
}
