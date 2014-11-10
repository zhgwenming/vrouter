package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
)

type cli struct {
	etcdClient *etcd.Client
}

func (c *cli) init(cmd *cobra.Command) {
}

func (c *cli) service(cmd *cobra.Command) {
}
