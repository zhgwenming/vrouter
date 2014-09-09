package daemon

import (
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/docker/docker/pkg/parsers/kernel"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/docker/libcontainer/netlink"
	"github.com/zhgwenming/vrouter/netinfo"
	"github.com/zhgwenming/vrouter/registry"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Daemon struct {
	etcdClient *etcd.Client

	// host relate information
	Hostname string

	// bridge information
	bridgeName  string
	bridgeIPNet *net.IPNet

	// interface information
	iface *net.Interface
	//ifaceIPNet *net.IPNet
	hostip string
}

func NewDaemon() *Daemon {
	return &Daemon{}
}

func (d *Daemon) listRoute() ([]*Route, uint64, error) {
	var index uint64
	var err error

	routes := make([]*Route, 0, 256)
	client := d.etcdClient

	routerPath := registry.RouterRoutesPrefix()
	if resp, err := client.Get(routerPath, false, true); err != nil {
		return routes, index, err
	} else {
		index = resp.EtcdIndex
		hosts := resp.Node.Nodes
		for _, host := range hosts {
			if hostKey := host.Key; strings.HasSuffix(hostKey, d.Hostname) {
				continue
			}
			value := host.Value
			r := ParseRoute(value)
			routes = append(routes, r)
		}
	}

	return routes, index, err
}

func (d *Daemon) ManageRoute() error {
	routes, etcdindex, err := d.listRoute()
	if err != nil {
		return err
	}

	for _, r := range routes {
		err = r.AddRoute(d.iface)
		if err != nil {
			log.Printf("error to add route: %s", err)
		}
	}

	receiver := make(chan *etcd.Response, 4)
	client := d.etcdClient

	go client.Watch(registry.RouterRoutesPrefix(), etcdindex, true, receiver, nil)

	//log.Printf("Watching or %s", registry.RouterHostsPrefix())

	for resp := range receiver {
		log.Printf("%v", resp.Node.Key)
	}

	return nil
}

func (d *Daemon) doKeepAlive(key, value string, ttl uint64) error {
	client := d.etcdClient

	if resp, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	} else {
		//log.Printf("No instance exist on this node, starting")
		go func() {
			sleeptime := time.Duration(ttl / 3)
			for {
				index := resp.EtcdIndex
				time.Sleep(sleeptime * time.Second)
				resp, err = client.CompareAndSwap(key, value, ttl, value, index)
				if err != nil {
					log.Fatal("Unexpected lost our node lock", err)
				}
			}
		}()
		return nil
	}
}

func (d *Daemon) KeepAlive() error {
	var err error
	if len(d.Hostname) == 0 {
		d.Hostname, err = os.Hostname()
		if err != nil {
			return err
		}
	}

	key := registry.NodeActivePath(d.Hostname)
	value := "alive"
	ttl := uint64(5)
	return d.doKeepAlive(key, value, ttl)
}

// return ip, ipnet, err
func (d *Daemon) getBridgeIPNet() (*net.IPNet, error) {
	client := d.etcdClient
	key := registry.BridgeInfoPath(d.Hostname)

	if resp, err := client.Get(key, false, false); err != nil {
		return nil, err
	} else {
		value := resp.Node.Value
		if ip, ipnet, err := net.ParseCIDR(value); err != nil {
			fmt.Printf("%v\n", value)
			return nil, err
		} else {
			ipnet.IP = ip
			return ipnet, nil
		}
	}
}

func (d *Daemon) updateRouterInterfaceNetIP(ip string) error {
	client := d.etcdClient

	key := registry.IfaceInfoPath(d.Hostname)
	value := ip
	ttl := uint64(0)

	if resp, err := client.Get(key, false, false); err == nil {
		if r := resp.Node.Value; r == value {
			log.Printf("found exist brideginfo for node: %s", d.Hostname)
			return nil
		}
	}

	// ignore response
	if _, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	}

	return nil
}

func (d *Daemon) updateNodeRoute() error {
	dnet := d.bridgeIPNet.Network()
	ip := d.hostip
	r := NewRoute(dnet, ip)

	client := d.etcdClient
	key := registry.NodeRoutePath(d.Hostname)
	value := r.String()
	ttl := uint64(0)

	if _, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	}
	return nil
}

// associate to nic ip address to an allocated IPNet
func (d *Daemon) BindBridgeIPNet(ifaceip string) (*net.IPNet, error) {
	var err error
	var brnet *net.IPNet

	if d.Hostname == "" {
		d.Hostname, err = os.Hostname()
		if err != nil {
			return brnet, err
		}
	}

	if ifaceip == "" {
		ifaceip = netinfo.GetFirstIPAddr()
	}

	// get node IPNet info first
	if brnet, err = d.getBridgeIPNet(); err != nil {
		return brnet, err
	}
	d.bridgeIPNet = brnet

	d.iface = netinfo.InterfaceByIPNet(ifaceip)

	if err = d.updateRouterInterfaceNetIP(ifaceip); err != nil {
		return brnet, err
	}
	err = d.updateNodeRoute()

	return brnet, err
}

func (d *Daemon) createBridgeIface(ifaceAddr string) error {
	kv, err := kernel.GetKernelVersion()
	// only set the bridge's mac address if the kernel version is > 3.3
	// before that it was not supported
	setBridgeMacAddr := err == nil && (kv.Kernel >= 3 && kv.Major >= 3)
	err = netlink.CreateBridge(d.bridgeName, setBridgeMacAddr)
	if err != nil {
		return err
	}

	iface, err := net.InterfaceByName(d.bridgeName)
	if err != nil {
		return err
	}

	ipAddr, ipNet, err := net.ParseCIDR(ifaceAddr)
	if err != nil {
		return err
	}

	if netlink.NetworkLinkAddIp(iface, ipAddr, ipNet); err != nil {
		return fmt.Errorf("Unable to add private network: %s", err)
	}
	if err := netlink.NetworkLinkUp(iface); err != nil {
		return fmt.Errorf("Unable to start network bridge: %s", err)
	}

	return nil
}

func WritePid(pidfile string) error {
	var file *os.File

	if _, err := os.Stat(pidfile); os.IsNotExist(err) {
		if file, err = os.Create(pidfile); err != nil {
			return err
		}
	} else {
		if file, err = os.OpenFile(pidfile, os.O_RDWR, 0); err != nil {
			return err
		}
		pidstr := make([]byte, 8)

		n, err := file.Read(pidstr)
		if err != nil {
			return err
		}

		if n > 0 {
			pid, err := strconv.Atoi(string(pidstr[:n]))
			if err != nil {
				fmt.Printf("err: %s, overwriting pidfile", err)
			}

			process, _ := os.FindProcess(pid)
			if err = process.Signal(syscall.Signal(0)); err == nil {
				return fmt.Errorf("pid: %d is running", pid)
			} else {
				fmt.Printf("err: %s, cleanup pidfile", err)
			}

			if file, err = os.Create(pidfile); err != nil {
				return err
			}

		}

	}
	defer file.Close()

	pid := strconv.Itoa(os.Getpid())
	fmt.Fprintf(file, "%s", pid)
	return nil
}
