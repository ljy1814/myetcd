package main

import (
	"fmt"
	"net"
	"time"
)

func ping() {
	addr, errAddr := net.ResolveTCPAddr("tcp4", "localhost:1234")
	if errAddr != nil {
		fmt.Printf("Client addr error: %s\n", errAddr)
		return
	}

	conn, errConn := net.DialTCP("tcp", nil, addr)
	defer conn.Close()
	if errConn != nil {
		fmt.Printf("Client conn error: %s\n", errConn)
		return
	}

	var buff [64]byte

	bytesTickerRead, errTickerRead := conn.Read(buff[0:])
	if errTickerRead != nil {
		fmt.Printf("Client read error: %s\n", errTickerRead)
		return
	}
	fmt.Printf("We got %s from the server.\n", string(buff[:bytesTickerRead]))

	msg := "Ping"
	_, errWrite := conn.Write([]byte(msg))
	if errWrite != nil {
		fmt.Printf("Client write error: %s\n", errWrite)
		return
	}

	bytesPingRead, errPingRead := conn.Read(buff[0:])
	if errPingRead != nil {
		fmt.Printf("Client read error: %s\n", errPingRead)
		return
	}

	fmt.Printf("we got %s from the server.\n", string(buff[:bytesPingRead]))
}

func main() {
	spanwInterval := time.Millisecond * 3000

	for {
		go ping()
		<-time.After(spanwInterval)
	}
}
