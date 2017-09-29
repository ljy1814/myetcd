package main

import (
	"encoding/json"
	"io/ioutil"
)

type EConfig struct {
	serviceMap        map[string][]string //服务端口映射表
	servicePartitions map[string]map[uint64]string
	Registry          `json:"registry"`
	Proxy             `json:"proxy"`
	Metrics           `json:"metrics"`
}

type Registry struct {
	Servers                   []string `json:"servers"`
	HeartbeatInternalInSecond int      `json:"heartbeat_interval_in_second"`
	HeartbeatTimeoutRound     int      `json:"heartbeat_timeout_round"`
}

type Proxy struct {
	HttpProxyDialTimeoutInMillisecond int `json:"http_proxy_dial_timeout_in_millisecond"`
	HttpProxyIoTimeoutInSecond        int `json:"http_proxy_io_timeout_in_second"`
}

type Metrics struct {
	ReportIntervalInSecond int    `json:"report_interval_in_second"`
	ReportAgentAddress     string `json:"report_agent_address"`
}

func (ec *EConfig) Load(path string) {
	var err error
	var jconf []byte
	jconf, err = ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jconf, &ec)
	if err != nil {
		panic(err)
	}
}
