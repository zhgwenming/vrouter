package daemon

import (
	"errors"
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/docker/libcontainer/netlink"
	"net"
	"strings"
)

type Route struct {
	target  string
	gateway string
}

func NewRoute(target, gateway string) *Route {
	return &Route{target: target, gateway: gateway}
}

func ParseRoute(str string) (*Route, error) {
	r := strings.Split(str, ":")
	if len(r) != 2 {
		return nil, errors.New("Wrong route format")
	}
	return &Route{target: r[0], gateway: r[1]}, nil
}

func (r *Route) AddRoute(src string, iface *net.Interface) error {
	return netlink.AddRoute(r.target, src, r.gateway, iface.Name)
}

func (r *Route) String() string {
	return r.target + ":" + r.gateway
}
