package main

import (
	"io"
	"net"
)

func TCPServer() {
	listener, _ := net.Listen("tcp", ":15001")
	for {
		in, _ := listener.Accept()
		go func() {
			out, _ := net.Dial("tcp", ":16001")
			go io.Copy(in, out)
			go io.Copy(out, in)
		}()
	}
}
