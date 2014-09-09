package daemon

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/docker/libcontainer/netlink"
	"net"
)

type Route struct {
	target string
	gw     string
}

func NewRoute(target, gw string) *Route {
	return &Route{target: target, gw: gw}
}

func ParseRoute(str string) *Route {
	return &Route{}
}

func (r *Route) AddRoute(iface *net.Interface) error {
	return netlink.AddRoute(r.target, "", r.gw, iface.Name)
}

func (r *Route) String() string {
	return r.target + ":" + r.gw
}
