package registry

import (
	"fmt"
	"os"

	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/etcd/config"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/etcd/server"
)

func StartEtcd() {
	var config = config.New()
	if err := config.Load(os.Args[1:]); err != nil {
		fmt.Println(server.Usage() + "\n")
		fmt.Println(err.Error() + "\n")
		os.Exit(1)
	} else if config.ShowVersion {
		fmt.Println("etcd version", server.ReleaseVersion)
		os.Exit(0)
	} else if config.ShowHelp {
		fmt.Println(server.Usage() + "\n")
		os.Exit(0)
	}

	var etcd = etcd.New(config)
	etcd.Run()
}
