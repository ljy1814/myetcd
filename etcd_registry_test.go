package myetcd

import (
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")

	addrs := []string{"http://127.0.0.1:2379"}
	r := NewEtcdRegistry(addrs)
	color.Red("registry : %v\n", r)
	err := r.Start()
	if err != nil {
		t.Error(err)
	}
	domain := "cao.chao"
	name := "mengde"
	version := "v1"
	addr := "localhost:80"

	err = r.DeleteService(domain, name, version)
	service := NewService(domain, name, version)
	color.Green("service : %v\n", service)
	service.NewEndpoint(addr)

	err = r.CreateService(service)
	ss, _ := r.GetService(domain, name, version)
	if ss != nil {
		color.Red("service : %v\n", ss)
	}
	if err != nil {
		t.Error(err)
	}
	color.Red("heartbeats : %v\n", r.heartbeats[servicePath(domain, name, version)])
	service, err = r.RegisterEndpoint(domain, name, version, addr, false)
	color.Yellow("service : %v, err[%v]\n", service, err)
	ep, _ := r.GetEndpoint(domain, name, version)
	color.Red("ep :%v\n", ep)
	time.Sleep(time.Second)
	color.Blue("heartbeat : %v\n", r.heartbeats)
	heartbeat := r.heartbeats[servicePath(domain, name, version)][0]
	if heartbeat.Addr != addr {
		t.Fatal(heartbeat)
	}

	r.Stop()
}
