package registry

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
	return VRouterPrefix() + "/" + node + "/" + "route"
}
