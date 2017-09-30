package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type EClient struct {
	httpClient *http.Client
}

func NewEClient() *EClient {
	return &EClient{httpClient: &http.Client{}}
}

func (ec *EClient) SendRequest(method string, oUrl string, oBody string) (result []byte) {
	request, err := http.NewRequest(method, oUrl, strings.NewReader(oBody))
	if err != nil {
		return
	}
	if len(oBody) > 0 {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}
	response, err := ec.httpClient.Do(request)
	if err != nil {
		return
	}
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	result, err = ioutil.ReadAll(response.Body)

	return
}
