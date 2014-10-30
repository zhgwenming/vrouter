package main

import (
	//"fmt"
	"github.com/zhgwenming/vrouter/controller"
	"github.com/zhgwenming/vrouter/daemon"
	"github.com/zhgwenming/vrouter/registry"
	"log"
)

func main() {

	// etcd client
	etcdConfig := new(registry.ClientConfig)

	// daemon instance
	cmd := daemon.NewConfig()

	routerCmd := cmd.InitCmd(etcdConfig)
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.Servers, "etcd_servers", "e", "https://127.0.0.1:4001", "etcd server uri")

	// cafile/certfile/keyfile
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.CaFile, "ca_file", "a", "", "etcd server ca file")
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.CertFile, "cert_file", "t", "", "etcd server cert file")
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.KeyFile, "key_file", "k", "", "etcd server key file")

	controller.InitCmd(routerCmd, etcdConfig)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
