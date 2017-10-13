package myetcd

import (
	"math/rand"
	"sync"
)

/*
 * Author : arch
 * Email : yajin160305@gmail.com
 * File : load_balance.go
 * CreateDate : 2017-10-12 17:04:31
 * */

/* vim: set tabstop=4 set shiftwidth=4 */

type LoadBalancer interface {
	GetEndpoint(path string, addrs []string) string
}

type RRLoadBalancer struct {
	indexes map[string]uint64
	rwMutex sync.RWMutex
}

func NewRRLoadBalancer() *RRLoadBalancer {
	return &RRLoadBalancer{indexes: make(map[string]uint64)}
}

// round robin load balancer
func (rr *RRLoadBalancer) GetEndpoint(path string, endpoints []string) (endpoint string) {
	rr.rwMutex.Lock()
	defer rr.rwMutex.Unlock()

	index := rr.indexes[path]
	pos := index % uint64(len(endpoints))
	endpoint = endpoints[pos]
	rr.indexes[path] = index + uint64(1)

	return endpoint
}

type RandomLoadBalancer struct {
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

func (rl *RandomLoadBalancer) GetEndpoint(path string, endpoints []string) (endpoint string) {
	if len(endpoints) <= 0 {
		return ""
	}

	index := rand.Intn(len(endpoints))
	return endpoints[index]
}

func NewLoadBalancer(policy string) LoadBalancer {
	switch policy {
	case "random":
		return NewRandomLoadBalancer()
	default:
		return NewRRLoadBalancer()
	}
}
