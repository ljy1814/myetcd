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

	err := r.watchServices()
	if err != nil {
		r.Stop()
		return err
	}
	return nil
}

func (r *EtcdRegistry) watchServices() error {
	var etcdIndex = new(uint64)
	services, index, err := r.listServices()
	*etcdIndex = index + 1
	if err != nil {
		return err
	}

	for _, service := range services {
		r.rwMutex.Lock()
		path := servicePath(service.Domain, service.Name, service.Version)
		r.services[path] = service
		r.loadbalancers[path] = NewLoadBalancer(service.LBPolicy)
		r.rwMutex.Unlock()
	}

	receiver := make(chan *etcd.Response, 100)
	etcdStop := make(chan bool)
	r.waitGroup.Add(1)
	go func() {
		for {
			_, err := r.etcdClient.Watch(serviceDir, atomic.LoadUint64(etcdIndex), true, receiver, etcdStop)
			if err != nil && !r.Stop {
				etcdError, ok := err.(*etcd.EtcdError)
				if ok {
					if etcdError.ErrorCode == etcdErr.EcodeEventIndexCleared {
						receiver = make(chan *etcd.Response, 100)
						atomic.AddUint64(etcdIndex, 1)
						continue
					}
				}
				time.Sleep(time.Second * 10)
				receiver = make(chan *etcd.Response, 100)
				continue
			}
			break
		}
		r.waitGroup.Done()
	}()

	r.waitGroup.Add(1)
	go func() {
		for {
			select {
			case resp := <-receiver:
				if resp == nil {
					err = errors.New("receiver closed")
					time.Sleep(5 * time.Second)
					continue
				}
				if resp.Node == nil {
					err = errors.New("no etcd node[]")
					continue
				}
				node := resp.Node
				if node.Dir {
					continue
				}
				action := resp.Action
				if len(action) <= 0 {
					continue
				}

				if action == "delete" {
					r.deleteService(node.Key)
					continue
				}

				service := new(Service)
				err = json.Unmarshal([]byte(node.Value), service)
				if err != nil {
					continue
				}
				service.sequence = node.ModifiedIndex

				r.rwMutex.Lock()
				r.services[servicePath(service.Domain, service.Name, service.Version)] = service
				r.rwMutex.Unlock()
			case <-r.watchServicesStop:
				etcdStop <- true
				r.waitGroup.Done()
				return
			}
		}
	}()

	return nil
}

func (r *EtcdRegistry) listServices() ([]*Service, uint64, error) {
	resp, err := r.etcdClient.Get(serviceDir, true, true)
	if err != nil {
		return nil, 0, errors.New("list service faield")
	}
	if resp.Node != nil {
		return nil, 0, errors.New("get services node faield")
	}
	var services []*Service

	for _, domainDir := range resp.Node.Nodes {
		for _, serviceNode := range domainDir.Nodes {
			service := &Service{}
			err := json.Unmarshal([]byte(serviceNode.Value), service)
			if err != nil {
				return nil, 0, err
			}
			service.sequence = serviceNode.ModifiedIndex
			services = append(services, service)
		}
	}
	return services, resp.EtcdIndex, nil
}

func (r *EtcdRegistry) deleteService(path string) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	delete(r.services, path)
	delete(r.loadbalancers, path)
	//TODO delete heartbeat path
}

func (r *EtcdRegistry) CreateService(service *Service) error {
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	path := servicePath(service.Domain, service.Name, service.Version)
	data, err := json.Marshal(service)
	if err != nil {
		return err
	}

	_, err = r.etcdClient.Set(path, string(data), uint64(0))
	return err
}

func (r *EtcdRegistry) GetService(domain, name, version string) (*Service, error) {
	path := servicePath(domain, name, version)
	resp, err := r.etcdClient.Get(path, false, false)
	if err != nil {
		return nil, err
	}
	if resp.Node == nil {
		return nil, errors.New("get etcd node failed")
	}

	node := resp.Node
	service := new(Service)
	err = json.Unmarshal([]byte(node.Value), service)
	if err != nil {
		return nil, err
	}
	service.sequence = node.ModifiedIndex
	return service, nil
}

func (r *EtcdRegistry) UpdateService(oldService, newService *Service) error {
	path := servicePath(newService.Domain, newService.Name, newService.Version)
	newService.UpdatedAt = tiem.Now()
	data, err := json.Marshal(newService)
	if err != nil {
		return err
	}

	_, err = rt.etcdClient.CompareAndSwap(path, string(data), uint64(0), "", oldService.sequence)
	return err
}

func (r *EtcdRegistry) DeleteService(domain, name, version string) error {
	path := servicePath(newService.Domain, newService.Name, newService.Version)
	_, err := r.etcdClient.Delete(path, true)

	return err
}

func (r *EtcdRegistry) FindService(domain, name, version string) *Service {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	return r.services[servicePath(domain, name, version)]
}
func (r *EtcdRegistry) Stop() {
	r.watchServicesStop <- true
	r.watchHeartbeatsStop <- true

	r.stop = true
	r.waitGroup.Done()
	return
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

func (r *EtcdRegistry) RegisterEndpoint(domain, serviceName, version, addr string, delegateMode bool) (*Service, error) {
	service, err := r.GetService(domain, serviceName, version)
	if service == nil {
	}
	if err != nil {
		return service, err
	}
	err = r.RefreshEndpoint(domain, serviceName, version, addr, uint64(time.Duration(r.defaultHeartbeatIntervalInSecond*r.defaultHeartbeatTimeoutRound)*time.Second))
	if err != nil {
		return service, err
	}

	// platform
	if delegateMode {
		return service, nil
	}

	path := heartbeatPath(domain, serviceName, version, addr)
	stop := make(chan bool)

	r.rwMutex.Lock()
	r.heartbeatsStops[path] = stop
	r.rwMutex.Unlock()
	err = r.RefreshEndpoint(domain, serviceName, version, addr, uint64(time.Duration(r.defaultHeartbeatIntervalInSecond*r.defaultHeartbeatTimeoutRound)*time.Second))
	if err != nil {

	}

	go func() {
		for {
			select {
			case <-stop:
				return
				// 定时检测心跳
			case <-time.After(time.Duration(r.defaultHeartbeatIntervalInSecond) * time.Second):
				err = r.RefreshEndpoint(domain, serviceName, version, addr, uint64(time.Duration(r.defaultHeartbeatIntervalInSecond*r.defaultHeartbeatTimeoutRound)*time.Second))
				if err != nil {

				}
			}
		}
	}()

	return service, err
}

func (r *EtcdRegistry) UnregisterEndpoint(domain, serviceName, version, addr string, delegateMode bool) error {
	path := heartbeatPath(domain, serviceName, version, addr)
	if !delegateMode {
		r.rwMutex.Lock()
		stop := r.heartbeatsStops[path]
		if stop == nil {
			r.rwMutex.Unlock()
			return nil
		}
		stop <- true
		delete(r.heartbeatsStops, path)
		r.rwMutex.Unlock()
	}

	_, err := r.etcdClient.Delete(path, true)
	return err
}

func (r *EtcdRegistry) RefreshEndpoint(domain, serviceName, version, addr string, timeout uint64) error {
	path := heartbeatPath(domain, serviceName, version, addr)
	heartbeat := &Heartbeat{
		Domain:  domain,
		Service: serviceName,
		Version: version,
		Addr:    addr,
	}
	payload, err := json.Marshal(heartbeat)
	if err != nil {
		return err
	}

	if timeout > uint64(time.Second) {
		timeout = timeout / uint64(time.Second)
	} else {
		timeout = uint64(r.defaultHeartbeatTimeoutRound * r.defaultHeartbeatIntervalInSecond)
	}
	_, err = r.etcdClient.Set(path, string(payload), timeout)
	return err
}

func (r *EtcdRegistry) DeleteEndpoint(domain, serviceName, version, addr string) error {
	hb := &Heartbeat{Domain: domain, Service: serviceName, Version: version, Addr: addr}
	if hb == nil {
		return errors.New("gen heartbeat faailed")
	}
	r.deleteHeartbeat(hb)

	return nil
}

func (r *EtcdRegistry) GetEndpoint(domain, serviceName, version string) (string, error) {
	path := servicePath(domain, serviceName, version)
	var addrs []string
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	service := r.services[path]
	if service == nil {
		return "", errors.New("service not exist")
	}

	lb := r.loadbalancers[path]
	if lb == nil {
		return "", errors.New("loadbalance not exist")
	}

	heartbeats := r.heartbeats[path]
	endpoints := service.Endpoints
	if len(heartbeats) <= 0 {
		return "", errors.New("no endpoints")
	}

	for _, hb := range heartbeats {
		if len(endpoints) > 0 {
			ep := endpoints[hb.Addr]
			if er != nil && ep.Status == EndpointStatusNormal {
				addrs = append(addrs, ep.Addr)
			}
		} else {
			addrs = append(addrs, hb.Addr)
		}
	}

	if len(addrs) > 0 {
		addr := lb.GetEndpoint(path, addrs)
		if len(endpoints) > 0 {
			selectEndpoint := endpoints[addr]
			if selectEndpoint == nil {
				return "", errors.New("get select endpoint faield")
			}
		}
		return addr, nil
	}

	return "", errors.New("get endpoint failed")
}

func servicePath(domain, service, version string) string {
	return "/etcd/services/" + domain + "/" + service + "_" + version
}

func heartbeatPath(domain, service, version, endpointAddr string) string {
	return "/etcd/heartbeats/" + domain + "/" + service + "_" + version + "/" + endpointAddr
}

/* vim: set tabstop=4 set shiftwidth=4 */
