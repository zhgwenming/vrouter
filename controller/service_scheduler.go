package controller

import (
	"github.com/zhgwenming/vrouter/service"
)

type ServiceScheduler struct {
	services    []*service.Service
	reScheduled []*service.Service
}
