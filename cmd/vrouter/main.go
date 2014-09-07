package main

import (
	//"fmt"
	"github.com/zhgwenming/vrouter/controller"
	"github.com/zhgwenming/vrouter/daemon"
	"log"
)

var (
	etcdServers string
)

func main() {

	cmd := daemon.NewCommand()
	routerCmd := cmd.InitCmd(&etcdServers)

	routerCmd.PersistentFlags().StringVarP(&etcdServers, "etcd_servers", "e", "http://127.0.0.1:4001", "etcd server uri")

	controller.InitCmd(routerCmd, &etcdServers)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
