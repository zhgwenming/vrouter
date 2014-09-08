package daemon

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/docker/libcontainer/netlink"
	"net"
)

type Route struct {
	bridgeIfaceAddr string
	routerIfaceAddr string
}

func (r *Route) AddRoute(iface *net.Interface) error {
	_, dnet, err := net.ParseCIDR(r.bridgeIfaceAddr)
	if err != nil {
		return err
	}

	ip, _, err := net.ParseCIDR(r.routerIfaceAddr)
	if err != nil {
		return err
	}

	err = netlink.AddRoute(dnet.String(), "", ip.String(), iface.Name)

	return err
}
