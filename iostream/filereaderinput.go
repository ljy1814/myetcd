package main

import (
	"fmt"
	"io"
	"os"
)

type alphaReader struct {
	reader io.Reader
}

func newAlphaReader(reader io.Reader) *alphaReader {
	return &alphaReader{
		reader: reader,
	}
}

func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func (a *alphaReader) Read(p []byte) (int, error) {
	n, err := a.reader.Read(p)
	if err != nil {
		return n, err
	}
	buf := make([]byte, n)
	j := 0
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			buf[j] = char
			j++
		}
	}
	copy(p, buf)
	return n, nil
}

func main() {
	fileReader, err := os.Open(("./filereaderinput.go"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fileReader.Close()
	reader := newAlphaReader(fileReader)
	p := make([]byte, 6)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}
