package registry


//v2/keys/_vrouter
//     ├── routes
//     │     ├── cluster1
//     │     │     ├── leader ── id {ip}
//     │     │     ├── resource
//     │     │     │    ├── createdIndex
//     │     │     │    ├── ...
//     │     │     │    └── createdIndexN
//     │     │     ├── node
//     │     │     │    ├── (node1) ip1 {pid}
//     │     │     │    ├── (node2) ip2
//     │     │     │    ├── ...
//     │     │     │    └── (nodeN) ipN
//     │     │     └── config
//     │     │
//     │     ├── clusterN
//     │
//     ├── hosts

const (
	DEFAULT_SUBNET  = "10.0.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

func VRouterPrefix() string {
	return REGISTRY_PREFIX + "/" + "route"
}

func RouterInterfacePath(node string) string {
	return VRouterPrefix() + "/" + node + "/" + "routerif"
}

func DockerBridgePath(node string) string {
	return VRouterPrefix() + "/" + node + "/" + "dockerbr"
}

func NodeRoutePath(node string) string {
	return REGISTRY_PREFIX + "/" + routes + "/" + node
}
