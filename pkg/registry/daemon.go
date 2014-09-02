package registry

import (
	"log"
	"time"
)

func (r *etcdRegistry) doKeepAlive(key, value string, ttl uint64) error {
	client := r.etcdClient

	if resp, err := client.Create(key, value, ttl); err != nil {
		log.Printf("Error to create node: %s", err)
		return err
	} else {
		//log.Printf("No instance exist on this node, starting")
		go func() {
			sleeptime := time.Duration(ttl / 3)
			for {
				index := resp.EtcdIndex
				time.Sleep(sleeptime * time.Second)
				resp, err = client.CompareAndSwap(key, value, ttl, value, index)
				if err != nil {
					log.Fatal("Unexpected lost our node lock", err)
				}
			}
		}()
		return nil
	}
}
