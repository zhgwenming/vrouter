package daemon

import (
	"fmt"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/netinfo"
	"github.com/zhgwenming/vrouter/registry"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

var (
	vrouter    *Daemon
	daemonMode bool
	gateway    bool
	hostname   string
	hostip     net.IP
)

type Daemon struct {
	etcdClient *etcd.Client
}

func NewDaemon(etcdClient *etcd.Client) *Daemon {
	return &Daemon{etcdClient: etcdClient}
}

func InitCmd() *cobra.Command {
	routerCmd := &cobra.Command{
		Use:  "vrouter",
		Long: "vrouter is a tool for routing distributed Docker containers.\n\n",
		Run:  Run,
	}

	// vrouter flags
	cmdflags := routerCmd.Flags()

	cmdflags.BoolVarP(&daemonMode, "daemon", "d", false, "whether to run as daemon mode")
	cmdflags.BoolVarP(&gateway, "gateway", "g", false, "to run as dedicated gateway, will not allocate subnet on this machine")
	cmdflags.StringVarP(&hostname, "hostname", "n", "", "hostname to use in daemon mode")
	cmdflags.IPVarP(&hostip, "hostip", "i", []byte{}, "use specified ip instead auto detected ip address")

	return routerCmd
}

func Run(c *cobra.Command, args []string) {
	if daemonMode {
		vrouter.KeepAlive(hostname)
		ipnet, err := BindHostNet(hostname, hostip)
		if err != nil {
			log.Fatal("Failed to bind ipnet, not initialized? - ", err)
		} else {
			log.Printf("main: get ipnet %v\n", ipnet)
		}
	} else {
		c.Help()
	}
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

func KeepAlive(hostname string) error {
	return vrouter.KeepAlive(hostname)
}

func (d *Daemon) getIPNet(hostname string) (*net.IPNet, error) {
	client := d.etcdClient
	key := registry.TenantNetPath(hostname)

	if resp, err := client.Get(key, false, false); err != nil {
		return nil, err
	} else {
		value := resp.Node.Value
		if _, ipnet, err := net.ParseCIDR(value); err != nil {
			fmt.Printf("%v\n", value)
			return nil, err
		} else {
			return ipnet, nil
		}
	}
}

func (d *Daemon) updateHostIP(hostname, ip string) error {
	client := d.etcdClient

	key := registry.HostNetPath(hostname)
	value := ip
	ttl := uint64(0)

	// ignore response
	if _, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	}

	return nil
}

// associate to nic ip address to an allocated IPNet
func (d *Daemon) BindHostNet(hostname, ip string) (*net.IPNet, error) {
	var err error
	var ipnet *net.IPNet

	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return ipnet, err
		}
	}

	if ip == "" {
		ip = netinfo.GetFirstIPAddr()
	}

	// get node IPNet info first
	if ipnet, err = d.getIPNet(hostname); err != nil {
		return ipnet, err
	}

	err = d.updateHostIP(hostname, ip)

	return ipnet, err
}

func BindHostNet(hostname string, ip net.IP) (*net.IPNet, error) {
	return vrouter.BindHostNet(hostname, string(ip))
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
