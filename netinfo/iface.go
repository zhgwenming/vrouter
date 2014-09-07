package netinfo

import (
	"net"
)

func InterfaceByIP(ip string) *net.Interface {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if addrs, err := i.Addrs(); err != nil {
			continue
		} else {
			for _, ipaddr := range addrs {
				//log.Printf("%v", ipaddr)
				ipnet, ok := ipaddr.(*net.IPNet)
				if !ok {
					continue
				}

				if ip == string(ipnet.IP) {
					return &i
				}
			}
		}
	}
	return nil
}
