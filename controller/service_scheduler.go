package controller

import (
	"github.com/zhgwenming/vrouter/service"
	"time"
)

type ServiceScheduler struct {
	run             chan struct{}
	lastRunTime     time.Time
	stableServices  map[string]*service.Service
	pendingServices []*service.Service
	orphanServices  []*service.Service
}

func NewServiceScheduler() *ServiceScheduler {
	sch := new(ServiceScheduler)
	run := make(chan struct{}, 1024)
	stable := make(map[string]*service.Service, 512)

	sch.run = run
	sch.stableServices = stable

	return sch
}

// Notify poke the scheduler as a event happened
func (srv *ServiceScheduler) Notify() {
	srv.run <- struct{}{}
}

// wait drain all the channel msg in case multiple events happened
func (srv *ServiceScheduler) wait() {
	<-srv.run
	for {
		_, ok := <-srv.run
		if !ok {
			break
		}
	}
}
