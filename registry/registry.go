package registry

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"strings"
)

type Registry struct {
	etcdClient *etcd.Client
}

const (
	DEFAULT_SUBNET  = "172.16.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

func RouterHostsPrefix() string {
	return REGISTRY_PREFIX + "/" + "hosts"
}

func RouterRoutesPrefix() string {
	return REGISTRY_PREFIX + "/" + "routes"
}

func NodeActivePath(node string) string {
	return REGISTRY_PREFIX + "/" + "members" + node
}

func RouterOverlayPath() string {
	return REGISTRY_PREFIX + "/" + "ipnet_overlay"
}

func IfaceInfoPath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "ifaceinfo"
}

func BridgeInfoPath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "bridgeinfo"
}

func NodeRoutePath(node string) string {
	return RouterRoutesPrefix() + "/" + node
}

func NewRegistry(etcd *etcd.Client) *Registry {
	return &Registry{etcd}
}

func (r *Registry) Create(key, value string, ttl uint64) error {
	client := r.etcdClient

	if _, err := client.Create(key, value, ttl); err != nil {
		//log.Printf("Error to create node: %s", err)
		return err
	}

	return nil
}

// set key to value
// create if not exist
func (r *Registry) Set(key, value string) error {
	client := r.etcdClient
	ttl := uint64(0)

	if resp, err := client.Get(key, false, false); err == nil {
		// exist, compare the value
		if resp.Node.Value != value {
			_, err = client.Update(key, value, ttl)
		}
		return err
	} else {
		if _, err = client.Create(key, value, ttl); err != nil {
			//log.Printf("Error to create node: %s", err)
			return err
		}
		return nil
	}
}

// List directory
func (r *Registry) List(prefix string) (map[string]string, uint64, error) {

	var index uint64
	var err error

	result := make(map[string]string, 256)
	client := r.etcdClient

	if resp, err := client.Get(prefix, true, true); err == nil {
		index = resp.EtcdIndex
		nodes := resp.Node.Nodes
		for _, node := range nodes {
			itemKey := node.Key
			itemKey = strings.TrimLeft(itemKey, prefix)
			value := node.Value

			result[itemKey] = value
		}
	}
	return result, index, err

}

// poll the specified prefix
func (r *Registry) Poll(prefix string, itemReceiver chan *Item) (map[string]string, error) {
	var err error
	var result map[string]string
	var index uint64

	client := r.etcdClient
	receiver := make(chan *etcd.Response, 4)

	if result, index, err = r.List(prefix); err == nil {
		go client.Watch(prefix, index, true, receiver, nil)
		go func() {
			//for resp := range receiver {
			//	node := resp.Node
			//	value := node.Value
			//	item = new(Item)
			//	item.key =

			//}
		}()
	}

	return result, err

}

func Set(client *etcd.Client, key, value string) error {
	ttl := uint64(0)

	if resp, err := client.Get(key, false, false); err == nil {
		// exist, compare the value
		if resp.Node.Value != value {
			_, err = client.Update(key, value, ttl)
		}
		return err
	} else {
		if _, err = client.Create(key, value, ttl); err != nil {
			//log.Printf("Error to create node: %s", err)
			return err
		}
		return nil
	}
}
