package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
	"reflect"
)

// InitCmd creates etcd client and register cobra subcommand
func InitCmd(parent *cobra.Command, client *registry.ClientConfig) {
	c := NewCli(client)

	v := reflect.ValueOf(c)
	for i := 0; i < v.NumMethod(); i++ {
		subcmd := v.Method(i).Interface().(func(*cobra.Command))
		subcmd(parent)
	}
}
