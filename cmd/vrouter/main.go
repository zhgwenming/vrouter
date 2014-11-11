package main

import (
	"github.com/zhgwenming/vrouter/controller"
	"github.com/zhgwenming/vrouter/daemon"
	"github.com/zhgwenming/vrouter/registry"
	"log"
	"os"
	"syscall"
)

func main() {

	tlsName := "tls"
	if hostname, err := os.Hostname(); err == nil {
		tlsName = hostname
	}

	// default tls cert info
	caFile := ""
	certFile := "/etc/vrouter/" + tlsName + ".crt"
	keyFile := "/etc/vrouter/" + tlsName + ".key"
	etcdServers := "https://127.0.0.1:4001"

	if file, found := syscall.Getenv("CA_FILE"); found {
		caFile = file
	}

	if file, found := syscall.Getenv("CERT_FILE"); found {
		certFile = file
	}

	if file, found := syscall.Getenv("KEY_FILE"); found {
		keyFile = file
	}

	if key, found := syscall.Getenv("ETCD_SERVERS"); found {
		etcdServers = key
	}

	if certFile == "" || keyFile == "" {
		etcdServers = "http://127.0.0.1:4001"
	}

	// etcd client
	etcdConfig := new(registry.ClientConfig)

	// daemon instance
	cmd := daemon.NewConfig()

	routerCmd := cmd.InitCmd(etcdConfig)
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.Servers, "etcd_servers", "e", etcdServers, "etcd server uri")

	// cafile/certfile/keyfile
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.CaFile, "ca-file", "", caFile, "etcd server ca file")
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.CertFile, "cert-file", "", certFile, "etcd server cert file")
	routerCmd.PersistentFlags().StringVarP(&etcdConfig.KeyFile, "key-file", "", keyFile, "etcd server key file")

	controller.AddSubCommands(routerCmd, etcdConfig)

	if err := routerCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
