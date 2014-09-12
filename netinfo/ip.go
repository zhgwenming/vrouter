package netinfo

import (
	"encoding/binary"
	"log"
	"net"
)

func ListIPNet(ip4Only bool) []*net.IPNet {
	ipnets := make([]*net.IPNet, 0, 4)
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}

		if addrs, err := i.Addrs(); err != nil {
			continue
		} else {
			for _, ipaddr := range addrs {
				//log.Printf("%v", ipaddr)
				ipnet, ok := ipaddr.(*net.IPNet)

				if !ok {
					log.Fatalf("assertion err: %v\n", ipnet)
				}

				ip4 := ipnet.IP.To4()
				if ip4Only && ip4 == nil {
					continue
				}

				if !ip4.IsLoopback() {
					ipnets = append(ipnets, ipnet)
				}
			}
		}
	}
	return ipnets
}

func GetFirstIPAddr() (addr string) {
	ifaces, _ := net.Interfaces()

iface:
	for _, i := range ifaces {
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}

		if addrs, err := i.Addrs(); err != nil {
			continue
		} else {
			for _, ipaddr := range addrs {
				//log.Printf("%v", ipaddr)
				ipnet, ok := ipaddr.(*net.IPNet)

				if !ok {
					log.Fatalf("assertion err: %v\n", ipnet)
				}

				ip4 := ipnet.IP.To4()
				if ip4 == nil {
					continue
				}
				//log.Printf("%v", ip4)

				if !ip4.IsLoopback() {
					addr = ipnet.String()
					break iface
				}
			}
		}
	}
	log.Printf("Found local ip4 %v", addr)
	return
}

func GetAllSubnet(ipnet *net.IPNet, hostbits int) []net.IPNet {
	ones, bits := ipnet.Mask.Size()
	zeros := bits - ones

	// network bits
	netBits := zeros - hostbits
	if netBits < 0 {
		return []net.IPNet{}
	}

	ip4 := ipnet.IP.To4()

	numberSubnet := 1 << uint(netBits)
	subnet := make([]net.IPNet, 0, numberSubnet)

	for i := uint32(0); i < uint32(numberSubnet); i++ {
		ipbuf := make([]byte, 4)
		number := i << uint(hostbits)
		binary.BigEndian.PutUint32(ipbuf, number)

		ip := (((uint32(ipbuf[0]) | uint32(ip4[0])) << 24) |
			((uint32(ipbuf[1]) | uint32(ip4[1])) << 16) |
			((uint32(ipbuf[2]) | uint32(ip4[2])) << 8) |
			uint32(ipbuf[3]) | uint32(ip4[3]+1))
		binary.BigEndian.PutUint32(ipbuf, ip)

		ipmask := net.CIDRMask(bits-hostbits, bits)

		subipnet := net.IPNet{IP: ipbuf, Mask: ipmask}
		subnet = append(subnet, subipnet)
	}

	return subnet

}
