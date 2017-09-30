package main

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	ec := &EConfig{}
	ec.Load("zc.json")
	fmt.Println(ec)
}
