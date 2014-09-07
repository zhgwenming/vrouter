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
	"syscall"
	"time"
)

type Daemon struct {
	etcdClient *etcd.Client
	iface      *net.Interface
	hostip     string
}

func NewDaemon(etcdClient *etcd.Client, ip string, iface *net.Interface) *Daemon {
	return &Daemon{etcdClient: etcdClient, hostip: ip, iface: iface}
}

func (d *Daemon) listRoute() ([]Route, uint64, error) {
	var index uint64
	var err error

	routes := make([]Route, 0, 256)
	client := d.etcdClient

	routerPath := registry.VRouterPrefix()
	if resp, err := client.Get(routerPath, false, true); err != nil {
		return routes, index, err
	} else {
		index = resp.EtcdIndex
		hosts := resp.Node.Nodes
		for _, host := range hosts {
			hostKey := host.Key

			// extract the host routerif/dockerbr info first
			routerNet := make(map[string]string, 2)
			for _, ipnet := range host.Nodes {
				routerNet[ipnet.Key] = ipnet.Value
			}

			routerIfacePath := hostKey + "/" + "routerif"
			routerIface := routerNet[routerIfacePath]
			if routerIface != "" && routerIface != d.hostip {
				dockerIfacePath := hostKey + "/" + "dockerbr"
				dockerIface := routerNet[dockerIfacePath]
				r := Route{routerIfaceAddr: routerIface, bridgeIfaceAddr: dockerIface}
				routes = append(routes, r)
			}
		}
	}

	return routes, index, err
}

func (d *Daemon) ManageRoute() error {
	routes, index, err := d.listRoute()
	if err != nil {
		return err
	}

	for r := range routes {
		r.AddRoute()
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

func (d *Daemon) KeepAlive(hostname string) error {
	var err error
	keyPrefix := registry.REGISTRY_PREFIX + "/" + "host"
	if len(hostname) == 0 {
		hostname, err = os.Hostname()
		if err != nil {
			return err
		}
	}

	key := keyPrefix + "/" + hostname
	value := "alive"
	ttl := uint64(5)
	return d.doKeepAlive(key, value, ttl)
}

func (d *Daemon) getDockerIPNet(hostname string) (*net.IPNet, error) {
	client := d.etcdClient
	key := registry.DockerBridgePath(hostname)

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

func (d *Daemon) updateRouterInterfaceNetIP(hostname, ip string) error {
	client := d.etcdClient

	key := registry.RouterInterfacePath(hostname)
	value := ip
	ttl := uint64(0)

	if resp, err := client.Get(key, false, false); err == nil {
		if r := resp.Node.Value; r == value {
			log.Printf("found exist routerif for node: %s", hostname)
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

// associate to nic ip address to an allocated IPNet
func (d *Daemon) BindDockerNet(hostname, ip string) (*net.IPNet, error) {
	var err error
	var hostnet *net.IPNet

	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return hostnet, err
		}
	}

	if ip == "" {
		ip = netinfo.GetFirstIPAddr()
	}

	// get node IPNet info first
	if hostnet, err = d.getDockerIPNet(hostname); err != nil {
		return hostnet, err
	}

	err = d.updateRouterInterfaceNetIP(hostname, ip)

	return hostnet, err
}

func (d *Daemon) createBridgeIface(bridgeIface, ifaceAddr string) error {
	kv, err := kernel.GetKernelVersion()
	// only set the bridge's mac address if the kernel version is > 3.3
	// before that it was not supported
	setBridgeMacAddr := err == nil && (kv.Kernel >= 3 && kv.Major >= 3)
	err = netlink.CreateBridge(bridgeIface, setBridgeMacAddr)
	if err != nil {
		return err
	}

	iface, err := net.InterfaceByName(bridgeIface)
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
