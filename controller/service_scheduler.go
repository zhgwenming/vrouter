package controller

import (
	"fmt"
	"github.com/zhgwenming/vrouter/service"
	"sync"
	"time"
)

type Notification struct {
	service uint64
	host    uint64
}

type ServiceScheduler struct {
	sync.Mutex
	run     chan struct{}
	lastRun time.Time

	// service related fields
	srvmap   map[string]*service.Service
	pendings []*service.Service
	orphans  []*service.Service

	// host related
	hosts []string
}

func NewServiceScheduler() *ServiceScheduler {
	sch := new(ServiceScheduler)
	run := make(chan struct{}, 1024)
	srvmap := make(map[string]*service.Service, 512)

	sch.run = run
	sch.srvmap = srvmap

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

	sch.srvmap[srv.Name] = srv
	sch.orphans = append(sch.orphans, srv)
	sch.notify()
}

func (sch *ServiceScheduler) DeleteService(srv *service.Service) {
	sch.Lock()
	defer sch.Unlock()

	delete(sch.srvmap, srv.Name)
	// not necessary to kick the scheduler
}

// add services on failed node to the orphan list
func (sch *ServiceScheduler) FailNode(node string) error {
	sch.Lock()
	defer sch.Unlock()

	var err = fmt.Errorf("Node not found %s", node)
	for _, srv := range sch.srvmap {
		if srv.Host == node {
			sch.orphans = append(sch.orphans, srv)
			err = nil
			break
		}
	}

	return err
}
