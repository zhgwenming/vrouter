package daemon

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/docker/libcontainer/netlink"
	"net"
)

type Route struct {
	bridgeIfaceAddr string
	routerIfaceAddr string
}

func (r *Route) AddRoute() error {
	_, dnet, err := net.ParseCIDR(r.bridgeIfaceAddr)
	if err != nil {
		return err
	}

	ip, ipnet, err := net.ParseCIDR(r.routerIfaceAddr)

	return nil
}
