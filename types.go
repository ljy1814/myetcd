package myetcd

import "time"

/*
 * Author : arch
 * Email : yajin160305@gmail.com
 * File : types.go
 * CreateDate : 2017-10-09 16:09:39
 * */

type Heartbeat struct {
	Domain  string `json:"domain"`
	Service string `json:"service"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
}

type Service struct {
	Domain          string               `json:"domain"`
	Name            string               `json:"name"`
	Version         string               `json:"version"`
	Type            string               `json:"type"`
	Status          string               `json:"status"`
	OnlyLeaderServe bool                 `json:"only_leader_serve,omitempty"`
	LBPolicy        string               `json:"lb_policy"`
	RetryTimes      uint8                `json:"retry_times"`
	DialTimeout     time.Duration        `json:"dial_timeout"`
	EndpointTimeout time.Duration        `json:"endpoint_timeout"`
	Endpoints       map[string]*Endpoint `json:"endpoints,omitempty"`
	Msgs            map[string]string    `json:"msgs"`
	CreatedAt       time.Time            `json:"created_at,omitempty"`
	UpdatedAt       time.Time            `json:"updated_at,omitempty"`
	sequence        uint64
}

type Endpoint struct {
	Addr string `json:"addr"`

	Status         string               `json:"status"`
	FreezeDuration time.Duration        `json:"freeze_duration,omitempty"`
	MsgFreezes     []*EndpointMsgFreeze `json:"msg_freezes,omitempty"`
}

const (
	EndpointStatusNormal  = "normal"
	EndpointStatusFreezed = "freezed"
)

/* vim: set tabstop=4 set shiftwidth=4 */
