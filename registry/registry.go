package registry

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"
)

type ServiceInfo struct {
	Name   string
	Host   string
	Port   int
	Weight int
}

type Registry struct {
	mu       sync.Mutex
	services map[string]ServiceInfo
}

func NewRegistry() *Registry {
	r := &Registry{
		services: make(map[string]ServiceInfo),
	}
	go r.startHeartbeat()
	return r
}

func (r *Registry) Register(info ServiceInfo, ok *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.services[info.Name] = info
	*ok = true
	return nil
}

func (r *Registry) Deregister(info ServiceInfo, ok *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.services, info.Name)
	*ok = true
	return nil
}

func (r *Registry) Lookup(_ struct{}, list *[]ServiceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	*list = (*list)[:0]
	for _, s := range r.services {
		*list = append(*list, s)
	}
	return nil
}

func (r *Registry) startHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		r.checkServices()
	}
}

func (r *Registry) checkServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, info := range r.services {
		addr := fmt.Sprintf("%s:%d", info.Host, info.Port)
		client, err := rpc.Dial("tcp", addr)
		if err != nil {
			log.Printf("Servizio %s non raggiungibile → rimosso", name)
			delete(r.services, name)
			continue
		}

		var ok bool
		err = client.Call("Service.HealthCheck", struct{}{}, &ok)
		client.Close()

		if err != nil || !ok {
			log.Printf("Servizio %s non risponde all'health-check → rimosso", name)
			delete(r.services, name)
		}
	}
}



