package registry

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
)

//v2/keys/_vrouter
//     ├── hosts
//     │     ├── hostname1
//     │     │     ├── active
//     │     │     ├── bridgeinfo
//     │     │     └── ifaceinfo
//     │     │
//     │     ├── ....
//     │     │
//     │     ├── hostnameN
//     │     │     ├── ...
//     │
//     ├── routes
//     │     ├── hostname1
//     │     ├── ...
//     │     ├── hostnameN

type Registry struct {
	etcdClient *etcd.Client
}

const (
	DEFAULT_SUBNET  = "10.0.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

func RouterHostsPrefix() string {
	return REGISTRY_PREFIX + "/" + "hosts"
}

func RouterRoutesPrefix() string {
	return REGISTRY_PREFIX + "/" + "routes"
}

func IfaceInfoPath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "ifaceinfo"
}

func BridgeInfoPath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "bridgeinfo"
}

func NodeActivePath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "active"
}

func NodeRoutePath(node string) string {
	return RouterRoutesPrefix() + "/" + node
}

func (r *Registry) Create(key, value string, ttl uint64) error {
	client := r.etcdClient

	if _, err := client.Create(key, value, ttl); err != nil {
		//log.Printf("Error to create node: %s", err)
		return err
	}

	return nil
}
