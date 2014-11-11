package controller

import (
	"github.com/zhgwenming/vrouter/service"
)

type ServiceScheduler struct {
	services    map[string]*service.Service
	reScheduled []*service.Service
}
