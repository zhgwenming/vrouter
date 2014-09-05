package main

import (
	//"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/controller"
	"github.com/zhgwenming/vrouter/daemon"
	"log"
	"strings"
)

var (
	etcdServer string
)

func main() {
	etcdstr := strings.Split(etcdServer, ",")
	etcdClient := etcd.NewClient(etcdstr)

	routerCmd := daemon.InitCmd()
	routerCmd.PersistentFlags().StringVarP(&etcdServer, "etcd_server", "e", "http://127.0.0.1:4001", "etcd daemon addr")

	controller.InitCmd(routerCmd, etcdClient)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
