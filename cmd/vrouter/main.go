package main

import (
	"github.com/zhgwenming/vrouter/controller"
	"github.com/zhgwenming/vrouter/daemon"
	"github.com/zhgwenming/vrouter/registry"
	"log"
	"os"
)

func main() {

	// etcd client
	etcdConfig := new(registry.ClientConfig)

	// daemon instance
	cmd := daemon.NewConfig()

	routerCmd := cmd.InitCmd(etcdConfig)
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.Servers, "etcd_servers", "e", "https://127.0.0.1:4001", "etcd server uri")

	tlsName := "tls"
	if hostname, err := os.Hostname(); err == nil {
		tlsName = hostname
	}

	// cafile/certfile/keyfile
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.CaFile, "ca-file", "", "", "etcd server ca file")
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.CertFile, "cert-file", "", "/etc/vrouter/"+tlsName+".crt", "etcd server cert file")
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.KeyFile, "key-file", "", "/etc/vrouter/"+tlsName+".key", "etcd server key file")

	controller.AddSubCommands(routerCmd, etcdConfig)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
