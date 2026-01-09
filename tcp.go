package main

import (
	"fmt"
	"io"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func TCPServer() {
	listener, _ := net.Listen("tcp", ":15001")
	for {
		in, _ := listener.Accept()
		go func() {
			sc, _ := in.(*net.TCPConn).SyscallConn()
			originalDst := &net.TCPAddr{}
			sc.Control(func(fd uintptr) {
				ipBytes, _ := syscall.GetsockoptInet4Addr(int(fd), syscall.IPPROTO_TCP, unix.SO_ORIGINAL_DST)
				originalDst := net.TCPAddr{
					IP:   net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3]),
					Port: 16001,
				}
				originalDst.IP = net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
				originalDst.Port = 16001
			})
			out, err := net.Dial("tcp", originalDst.String())
			if err != nil {
				defer in.Close()
				fmt.Printf("can not connect to original dst: %v\n", err)
				return
			}
			go io.Copy(in, out)
			go io.Copy(out, in)
		}()
	}
}
