package main

import (
	//"fmt"
	"github.com/zhgwenming/vrouter/controller"
	"github.com/zhgwenming/vrouter/daemon"
	"log"
)

var (
	etcdServer string
)

func main() {
	routerCmd := daemon.InitCmd()
	routerCmd.PersistentFlags().StringVarP(&etcdServer, "etcd_server", "e", "http://127.0.0.1:4001", "etcd daemon addr")

	controller.InitCmd(routerCmd, etcdServer)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
