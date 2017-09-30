package main

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestSendReq(t *testing.T) {
	ec := NewEClient()
	if ec == nil {
		fmt.Println("new ec client failed")
		return
	}

	method := "PUT"
	oUrl := "http://172.16.127.37:4001/v2/keys/songjiang"
	param := url.Values{}
	param.Set("value", "天罡星")
	//param.Set("value", "yi")
	param.Set("ttl", "5")
	resp := ec.SendRequest(method, oUrl, param.Encode())
	fmt.Printf("put resp : %v\n", string(resp))

	// 还未删除
	time.Sleep(5201 * time.Millisecond)

	method = "GET"
	resp = ec.SendRequest(method, oUrl, "")
	fmt.Printf("get resp : %v\n", string(resp))

	method = "DELETE"
	resp = ec.SendRequest(method, oUrl, "")
	fmt.Printf("delete resp : %v\n", string(resp))

	method = "GET"
	resp = ec.SendRequest(method, oUrl, "")
	fmt.Printf("get resp : %v\n", string(resp))

}
