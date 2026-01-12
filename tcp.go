package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const idleTimeout = 10 * time.Second

func TCPServer() {
	listener, _ := net.Listen("tcp", ":15001")
	for {
		in, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept error: %v\n", err)
			continue
		}
		go func() {
			defer in.Close()
			inConn := getTCPConn(in)
			originalDst, err := getOriginalDst(inConn)
			if err != nil {
				fmt.Printf("failed to obtain original dst: %v\n", err)
				return
			}

			out, err := net.DialTimeout("tcp", originalDst.String(), 5*time.Second)
			if err != nil {
				fmt.Printf("can not connect to original dst: %v, err: %v\n", originalDst.String(), err)
				return
			}
			defer out.Close()
			outConn := getTCPConn(out)
			wg := sync.WaitGroup{}
			transfer(inConn, outConn, &wg)
			transfer(outConn, inConn, &wg)
			wg.Wait()
		}()
	}
}
