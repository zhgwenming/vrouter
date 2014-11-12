package controller

import (
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
	"github.com/zhgwenming/vrouter/service"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type ServiceConfig struct {
	Name     string
	Listen   string
	BackEnds string
}

type ServiceManager struct {
	srvConfig ServiceConfig
	config    *Config

	etcdClient *etcd.Client
}

func (mgr *ServiceManager) Run(cmd *cobra.Command, args []string) {
	var action string

	mgr.etcdClient = registry.NewClient(mgr.config.etcdConfig)

	if len(args) > 0 {
		var err error
		action = args[0]
		// all the actions
		switch action {
		case "add":
			err = mgr.Add()
		case "delete":
			err = mgr.Delete()
		case "ls":
			fallthrough
		case "list":
			err = mgr.List()
		default:
			cmd.Usage()
		}

		if err != nil {
			log.Fatalf("Error to add service: %s\n", err)
		}
	} else {
		fmt.Printf("No action specified.\n")
		cmd.Usage()
	}

}

func (mgr *ServiceManager) Add() error {
	if mgr.srvConfig.Name == "" {
		return fmt.Errorf("No service name specified")
	}

	key := registry.RouterServicesPrefix() + "/" + mgr.srvConfig.Name

	srv := service.NewService()
	srv.Name = mgr.srvConfig.Name

	// listen addr/port
	listen := strings.Split(mgr.srvConfig.Listen, ":")
	if len(listen) != 2 {
		return fmt.Errorf("error format of listen: %s, should be ip:port form")
	}

	srv.Addr = listen[0]
	srv.Port = listen[1]
	srv.Active = true
	srv.CreateTime = time.Now()

	for _, backend := range strings.Fields(mgr.srvConfig.BackEnds) {
		if b, err := service.NewBackend(backend); err == nil {
			srv.AddBackend(b)
		} else {
			return err
		}
	}

	value := string(srv.Marshal())
	if _, err := mgr.etcdClient.Create(key, value, uint64(0)); err != nil {
		return err
	}
	fmt.Printf("service %s added\n", mgr.srvConfig.Name)

	return nil
}

func (mgr *ServiceManager) Delete() error {
	if mgr.srvConfig.Name == "" {
		return fmt.Errorf("No service name specified")
	}

	key := registry.RouterServicesPrefix() + "/" + mgr.srvConfig.Name
	if _, err := mgr.etcdClient.Delete(key, true); err != nil {
		return err
	}

	fmt.Printf("service %s deleted\n", mgr.srvConfig.Name)
	return nil
}

func (mgr *ServiceManager) List() error {
	//fmt.Printf("list services:\n")

	key := registry.RouterServicesPrefix()
	resp, err := mgr.etcdClient.Get(key, true, true)
	if err != nil {
		return err
	}

	nodes := resp.Node.Nodes

	length := len(nodes)
	services := make([]service.Service, length)
	for i, n := range nodes {
		value := n.Value

		//fmt.Printf("value is %#v", value)
		services[i].UnMarshal(value)
	}

	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "NAME\tADDRESS PORT\tCREATED TIME\tACTIVE\tBACKENDS\t")

	for _, s := range services {
		var status string
		if s.Active {
			status = "yes"
		} else {
			status = "no "
		}
		fmt.Fprintf(w, "%s\t%s:%s\t"+"%s\t"+"%s\t%s\n",
			s.Name, s.Addr, s.Port,
			s.CreateTime.Local().Format("2006-01-01 15:04:05"),
			status, s.Backends)
	}
	//fmt.Fprintln(w)
	w.Flush()

	return nil
}

// Schedule assign a active service a active host
func (mgr *ServiceManager) Schedule() {
	key := registry.RouterServicesPrefix()
	resp, err := mgr.etcdClient.Get(key, true, true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	nodes := resp.Node.Nodes
	//index := resp.EtcdIndex

	length := len(nodes)
	services := make([]service.Service, length)
	for i, n := range nodes {
		value := n.Value

		services[i].UnMarshal(value)
	}
}

// ServeNode would run the actuall service on all the slave nodes
func (mgr *ServiceManager) ServeNode(node string) {
}
