package registry

//v2/keys/_vrouter
//     ├── hosts
//     │     ├── hostname1
//     │     │     ├── active
//     │     │     ├── bridgeinfo
//     │     │     └── ifaceinfo
//     │     │
//     │     ├── ....
//     │     │
//     │     ├── hostnameN
//     │     │     ├── ...
//     │
//     ├── routes
//     │     ├── hostname1
//     │     ├── ...
//     │     ├── hostnameN

const (
	DEFAULT_SUBNET  = "10.0.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

func RouterHostsPrefix() string {
	return REGISTRY_PREFIX + "/" + "hosts"
}

func RouterRoutesPrefix() string {
	return REGISTRY_PREFIX + "/" + "routes"
}

func IfaceInfoPath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "ifaceinfo"
}

func BridgeInfoPath(node string) string {
	return RouterHostsPrefix() + "/" + node + "/" + "bridgeinfo"
}

func NodeRoutePath(node string) string {
	return RouterRoutesPrefix() + "/" + node
}
