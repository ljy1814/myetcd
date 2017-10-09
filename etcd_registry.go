package myetcd

import (
	"encoding/json"
	"errors"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	etcd "github.com/cores/go-etcd"
)

/*
 * Author : arch
 * Email : yajin160305@gmail.com
 * File : etcd_registry.go
 * CreateDate : 2017-10-09 11:54:46
 * */

var (
	heartbeatDir = "/etcd/heartbeats"
	serviceDir   = "/etcd/services"
)

type EtcdRegistry struct {
	// ETCD client
	etcdClient          *etcd.Client
	services            map[string]*Service
	watchServicesStop   chan bool
	watchHeartbeatsStop chan bool

	heartbeatsStops                  map[string]chan bool
	heartbeats                       map[string][]*Heartbeat
	loadbalancers                    map[string]LoadBalancer
	defaultHeartbeatIntervalInSecond int
	defaultHeartbeatTimeoutRound     int

	stop      bool
	waitGroup sync.WaitGroup
	rwMutex   sync.RWMutex

	etcdAddrs []string
}

func NewEtcdRegistry(addrs []string) *EtcdRegistry {
	cli := etcd.NewClient(addrs)
	return &EtcdRegistry{
		etcdClient:          cli,
		services:            make(map[string]*Service),
		watchServicesStop:   make(chan bool),
		watchHeartbeatsStop: make(chan bool),

		heartbeatsStops:                  make(map[string]chan bool),
		heartbeats:                       make(map[string][]*Heartbeat),
		loadbalancers:                    make(map[string]LoadBalancer),
		defaultHeartbeatIntervalInSecond: 10,
		defaultHeartbeatTimeoutRound:     2,
	}
}

func (r *EtcdRegistry) Sync() error {
	succ := r.etcdClient.SyncCluster()
	if !succ {
		return errors.New("cannot sync machines")
	}
	return nil
}

func (r EtcdRegistry) Start() error {
	err := r.Sync()
	if err != nil {
		panic(err)
	}

	go r.watchHeartbeats()
}

func (r *EtcdRegistry) watchHeartbetas() error {
loop:
	var etcdIndex = new(uint64)
	heartbeatList, index, err := r.listHeartbeats()
	*etcdIndex = index + 1
	if err != nil {
		panic(err)
	}
	heartbeats := make(map[string][]*Heartbeat)

	for _, heartbeat := range heartbeatList {
		path := servicePath(heartbeat.Domain, heartbeat.Service, heartbeat.Version)
		heartbeats[path] = append(heartbeats[path], heartbeat)
	}
	r.resetHeartbeats(heartbeats)

	for {
		receiver := make(chan *etcd.Response, 5000)
		watchDone := make(chan error)
		go func() {
			_, err := r.etcdClient.Watch(heartbeatDir, atomic.LoadUint64(etcdIndex), true, receiver, nil)
			if err != nil {
				watchDone <- err
				return
			}
		}()

		for resp := range receiver {
			if resp == nil {
				break
			}
			if resp.Node == nil {
				continue
			}
			node := resp.Node

			atomic.StoreUint64(etcdIndex, resp.Node.ModifiedIndex+1)
			action := resp.Action
			if len(resp.Action) <= 0 {
				break
			}

			if node.Dir {
				continue
			}

			if action == "delete" || action == "expire" {
				var heartbeat *Heartbeat
				keyRegexp := "^/etcd/heartbeats/(?P<domain>[^/]+)/(?P<serviceName>[^_]+)_(?P<version>[^/]+)/(?P<addr>.+)$"
				heartbeatKeyRegexp, err := regexp.Compile(keyRegexp)
				if err != nil {
					heartbeatKeyRegexp = nil
					return err
				}
				matches := heartbeatKeyRegexp.FindStringSubmatch(node.Key)
				if matches != nil && len(matches) == 5 {
					heartbeat = new(Heartbeat)
					heartbeat.Domain = matches[1]
					heartbeat.Service = matches[2]
					heartbeat.Version = matches[3]
					heartbeat.Addr = matches[4]
				} else {

				}
				if heartbeat == nil {
					heartbeat = new(Heartbeat)
					err = json.Unmarshal([]byte(resp.PrevNode.Value), heartbeat)
					if err != nil {
						heartbeat = nil
					}
				}
				if heartbeat != nil {
					r.deleteHeartbeat(heartbeat)
				}
			} else {
				heartbeat = new(Heartbeat)
				err := json.Unmarshal([]byte(node.Value), heartbeat)
				if err != nil {

				}
				r.addHeartbeat(heartbeat)
			}
		}
		watchErr := <-watchDone
		if e, ok := watchErr.(*etcd.EtcdError); ok && e.ErrorCode == 401 {
			time.Sleep(time.Second * 10)
			goto loop
		} else {
			time.Sleep(time.Second * 1)
		}

		select {
		case <-r.watchHeartbeatsStop:
			return nil
		default:
		}
	}
}

func (r *EtcdRegistry) listHeartbeats() ([]*Heartbeat, uint64, error) {
	resp, err := r.etcdClient.Get(heartbeatDir, true, true)
	if err != nil {
		return nil, 0, err
	}
	if resp.Node == nil {
		return nil, 0, errors.New("check etcd node error")
	}

	var heartbeats []*Heartbeat
	for _, domainDir := range resp.Node.Nodes {
		for _, serviceDir := range domainDir.Nodes {
			for _, heartbeatNode := range serviceDir.Nodes {
				heartbeat := &Heartbeat{}
				err = json.Unmarshal([]byte(heartbeatNode.Value), heartbeat)
				if err != nil {
					return nil, uint64(0), err
				}
			}
			heartbeats = append(heartbeats, heartbeat)
		}
	}
	return heartbeats, resp.EtcdIndex, nil
}

func (r *EtcdRegistry) deleteHeartbeat(heartbeat *Heartbeat) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	path := servicePath(heartbeat.Domain, heartbeat.Service, heartbeat.Version)

	var heartbeats []*Heartbeat

	for _, h := range r.heartbeats[path] {
		if h.Addr != heartbeat.Addr {
			heartbeats = append(heartbeats, h)
		}
	}
	if len(heartbeats) > 0 {
		r.heartbeats[path] = heartbeats
	} else {
		delete(r.heartbeats, path)
	}
}

func (r *EtcdRegistry) addHeartbeat(heartbeat *Heartbeat) {
	r.rwMutex.RLock()
	path := servicePath(heartbeat.Domain, heartbeat.Service, heartbeat.Version)
	for _, h := range r.heartbeats[path] {
		if h.Addr == heartbeat.Addr {
			r.rwMutex.RUnlock()
			return
		}
	}
	r.rwMutex.RUnlock()

	r.rwMutex.Lock()
	r.heartbeats[path] = append(r.heartbeats[path], heartbeat)
	r.rwMutex.Unlock()
}

func (r *EtcdRegistry) resetHeartbeats(heartbeats map[string][]*Heartbeat) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	r.heartbeats = heartbeats
}

func servicePath(domain, service, version string) string {
	return "/etcd/services/" + domain + "/" + service + "_" + version
}

/* vim: set tabstop=4 set shiftwidth=4 */
