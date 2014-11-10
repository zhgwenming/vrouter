package controller

type cli struct {
	etcdClient *etcd.Client
}

func (c *cli) init(cmd *cobra.Command) {
}

func (c *cli) service(cmd *cobra.Command) {
}
