package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("vim-go")
	fff()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	})
	http.ListenAndServe(":1999", nil)
}

func fff() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	go func() {
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}()
}
