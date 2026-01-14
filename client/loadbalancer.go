package main

import (
	"math/rand"
	"SDCC_GO/registry"
)

var rrIndex=0

//Caso stateless Random
func RandomLB(services []registry.ServiceInfo) registry.ServiceInfo {
    if len(services) == 0 {
        return registry.ServiceInfo{} // oppure panic controllato
    }
    return services[rand.Intn(len(services))]
}

//Caso stateless RoundRobin
func RoundRobinLB(services []registry.ServiceInfo) registry.ServiceInfo{
	s:=services[rrIndex%len(services)]
	rrIndex++
	return s
}

//Caso statefull pesato 
func WeightedLB(services []registry.ServiceInfo) registry.ServiceInfo{
	total:=0
	for _,s := range services{
		total += s.Weight
	}
	r:=rand.Intn(total)
	for _,s := range services {
		if r<s.Weight{
			return s
		}
		r-=s.Weight
	}
	return services[0]
}