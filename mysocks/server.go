package main

import (
	"container/list"
	"fmt"
	"net"
	"runtime"
	"time"
)

type Client struct {
	Name      string
	Conn      net.Conn
	ReadChan  chan string
	WriteChan chan string
	CloseChan chan bool
}

func NewClient(name string, conn net.Conn) *Client {
	fmt.Printf("%s joined.\n", name)
	newClient := &Client{
		name,
		conn,
		make(chan string),
		make(chan string),
		make(chan bool),
	}

	return newClient
}

func (c *Client) Close() {
	close(c.CloseChan)
	c.Conn.Close()
	fmt.Printf("%s left.\n", c.Name)
}

func (c *Client) ReadToReadChan() {
	defer c.Close()

	var buffer [1024]byte
	for {
		bytesRead, errRead := c.Conn.Read(buffer[0:])
		if errRead != nil {
			fmt.Printf("! Client read error: %s.\n", errRead)
			return
		}

		request := string(buffer[:bytesRead])
		c.ReadChan <- request
	}
}

func (c *Client) WriteChan2Write() {
	for {
		select {
		case <-c.CloseChan:
			return
		case response := <-c.WriteChan:
			fmt.Printf("Server => %s : %s\n", c.Name, response)
			_, errWrite := c.Conn.Write([]byte(response))
			if errWrite != nil {
				fmt.Printf("! Client write error: %s.\n", errWrite)
				return
			}
		}
	}
}

type Clients struct {
	list *list.List
}

func NewClients() *Clients {
	return new(Clients).Init()
}

func (c *Clients) Init() *Clients {
	c.list = list.New()
	return c
}

func (c *Clients) Count() int {
	return c.list.Len()
}

func (c *Clients) SendAll(msg string) {
	for entry := c.list.Front(); entry != nil; entry = entry.Next() {
		client := entry.Value.(*Client)
		client.WriteChan <- msg
	}
}

func (c *Clients) Add(client *Client) {
	c.list.PushBack(client)
}

func (c *Clients) RemoveOnClose(client *Client) {
	<-client.CloseChan //等待client关闭
	for entry := c.list.Front(); entry != nil; entry = entry.Next() {
		otherClient := entry.Value.(*Client)
		if otherClient.Name == client.Name {
			c.list.Remove(entry)
		}
	}
}

func main() {
	addr, _ := net.ResolveTCPAddr("tcp4", ":1234")
	listener, _ := net.ListenTCP("tcp", addr)
	var clients = NewClients()

	infoTicker := time.NewTicker(time.Second)
	go func() {
		for {
			<-infoTicker.C

			c := clients.Count()
			g := runtime.NumGoroutine()
			fmt.Printf("%d clients online, running on %d goroutines.\n", c, g)
		}
	}()

	tickTockTicker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			<-tickTockTicker.C
			clients.SendAll("Tick")
			<-tickTockTicker.C
			clients.SendAll("Tock")
		}
	}()

	clientNumber := 0
	for {
		conn, _ := listener.Accept()
		clientNumber++
		name := fmt.Sprintf("Client[%d]", clientNumber)
		client := NewClient(name, conn)
		clients.Add(client)

		go clients.RemoveOnClose(client)
		go handleClient(client)
	}
}

func handleClient(client *Client) {
	for {
		select {
		case <-client.CloseChan:
			return
		case request := <-client.ReadChan:
			fmt.Printf("%s => Server : %s\n", client.Name, request)
			client.WriteChan <- "Pong"
		}
	}
}
