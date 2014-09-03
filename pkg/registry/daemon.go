package registry

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

func (r *Registry) doKeepAlive(key, value string, ttl uint64) error {
	client := r.etcdClient

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

func (r *Registry) KeepAlive(hostname string) error {
	var err error
	keyPrefix := REGISTRY_PREFIX + "/" + "host"
	if len(hostname) == 0 {
		hostname, err = os.Hostname()
		if err != nil {
			return err
		}
	}

	key := keyPrefix + "/" + hostname
	value := "alive"
	ttl := uint64(5)
	return r.doKeepAlive(key, value, ttl)
}

func KeepAlive(hostname string) error {
	return registryClient.KeepAlive(hostname)
}

func (r *Registry) UpdateHostIP(hostname, ip string) error {
	var err error

	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return err
		}
	}

	if ip == "" {
		ip = GetFirstIPAddr()
	}

	client := r.etcdClient

	key := registryRoutePrefix() + "/" + hostname + "/" + "ipaddr"
	value := ip
	ttl := uint64(0)

	// ignore response
	if _, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	}

	return nil
}

func UpdateHostIP(hostname string, ip net.IP) error {
	return registryClient.UpdateHostIP(hostname, string(ip))
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
