package controller

import (
	"github.com/zhgwenming/vrouter/service"
	"sync"
	"time"
)

type ServiceScheduler struct {
	sync.Mutex
	run      chan struct{}
	lastRun  time.Time
	stables  map[string]*service.Service
	pendings []*service.Service
	orphans  []*service.Service
}

func NewServiceScheduler() *ServiceScheduler {
	sch := new(ServiceScheduler)
	run := make(chan struct{}, 1024)
	stable := make(map[string]*service.Service, 512)

	sch.run = run
	sch.stables = stable

	return sch
}

// Notify poke the scheduler as a event happened
func (sch *ServiceScheduler) notify() {
	sch.run <- struct{}{}
}

// wait drain all the channel msg in case multiple events happened
func (sch *ServiceScheduler) wait() {
	<-sch.run
	for {
		_, ok := <-sch.run
		if !ok {
			break
		}
	}
}

func (sch *ServiceScheduler) AddService(srv *service.Service) {
	sch.Lock()
	defer sch.Unlock()

	sch.orphans = append(sch.orphans, srv)
	sch.notify()
}

func (sch *ServiceScheduler) FailNode(node string) {
	for key, srv := range sch.stables {
		if srv.Host == node {
			sch.Lock()
			defer sch.Unlock()

			delete(sch.stables, key)
			sch.orphans = append(sch.orphans, srv)
			return
		}
	}
}
