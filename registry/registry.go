package registry

const (
	DEFAULT_SUBNET  = "10.0.0.0/16"
	REGISTRY_PREFIX = "_vrouter"
)

func RoutePrefix() string {
	return REGISTRY_PREFIX + "/" + "route"
}

func TenantNetPath(node string) string {
	return RoutePrefix() + "/" + node + "/" + "tenantnet"
}

func HostNetPath(node string) string {
	return RoutePrefix() + "/" + node + "/" + "hostnet"
}
