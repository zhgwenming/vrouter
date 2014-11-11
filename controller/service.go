package controller

import (
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/registry"
	"github.com/zhgwenming/vrouter/service"
	"log"
	"strings"
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
		case "list":
			err = mgr.List()
		default:
			cmd.Usage()
		}

		if err != nil {
			log.Fatalf("Error to add service: %s\n", err)
		}
		fmt.Printf("%#v\n", action)
	} else {
		// list all exist service?
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
	fmt.Printf("List services\n")

	return nil
}
