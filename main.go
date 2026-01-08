package main

import (
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: main <upstream-address>")
	}
	upstreamAddr := os.Args[1]

	listener, err := net.Listen("tcp", ":15001")
	if err != nil {
		panic(err)
	}

	for {
		in, _ := listener.Accept()
	}
}
