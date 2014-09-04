package main

import (
	//"fmt"
	"github.com/zhgwenming/vrouter/daemon"
	"log"
)

var (
	etcdServer string
)

func main() {
	routerCmd := daemon.DaemonInit()
	routerCmd.PersistentFlags().StringVarP(&etcdServer, "etcd_server", "e", "http://127.0.0.1:4001", "etcd daemon addr")

	daemon.Init(routerCmd, etcdServer)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
