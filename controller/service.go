package controller

import (
	"github.com/zhgwenming/vrouter/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zhgwenming/vrouter/service"
)

type Service struct {
	service.Service
	cmd *Command
}

func (srv *Service) Run(cmd *cobra.Command, args []string) {
}
