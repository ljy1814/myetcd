package main

import (
	"bytes"
	"io"
	"os"
)

func main() {
	proverbs := new(bytes.Buffer)
	proverbs.WriteString("channels or chestrate mutexes serialize\n")
	proverbs.WriteString("cgo is not go\n")
	proverbs.WriteString("errors are values\n")
	proverbs.WriteString("don't panic\n")

	piper, pipew := io.Pipe()
	go func() {
		defer pipew.Close()
		// 将数据写入管道
		io.Copy(pipew, proverbs)
	}()

	// 将数据从管道读出
	io.Copy(os.Stdout, piper)
	piper.Close()
}
