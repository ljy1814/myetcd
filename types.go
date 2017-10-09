package myetcd

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

/* vim: set tabstop=4 set shiftwidth=4 */
