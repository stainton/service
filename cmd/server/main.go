package main

import (
	"fmt"
	"net"
	"os"
)

func l4() {
	listener, err := net.Listen("tcp", ":15001")
	if err != nil {
		fmt.Printf("listen on 15001 failed with error: %v\n", err)
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept connection failed with error: %v\n", err)
			continue
		}
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			fmt.Printf("not a TCP connection\n")
			conn.Close()
			continue
		}
		// tcpConn.
		// buffer := make([]byte, 1024)
		// if l7Grpc(tcpConn, buffer) {
		// 	continue
		// } else if l7Http(tcpConn, buffer) {
		// 	continue
		// } else {
		// 	l4Transport(tcpConn, buffer)
		// }
	}
}
