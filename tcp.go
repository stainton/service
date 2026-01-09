package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

const idleTimeout = 10 * time.Second

func TCPServer() {
	listener, _ := net.Listen("tcp", ":15001")
	for {
		in, _ := listener.Accept()
		go func() {
			defer in.Close()
			_ = in.SetDeadline(time.Now().Add(idleTimeout))
			sc, _ := in.(*net.TCPConn).SyscallConn()
			originalDst := &net.TCPAddr{Port: 16001}
			var dstErr error
			sc.Control(func(fd uintptr) {
				// SO_ORIGINAL_DST returns a sockaddr_in in IPv6Mreq.Multiaddr
				addr, err := unix.GetsockoptIPv6Mreq(int(fd), unix.IPPROTO_IP, unix.SO_ORIGINAL_DST)
				if err != nil {
					dstErr = fmt.Errorf("getsockopt SO_ORIGINAL_DST: %w", err)
					return
				}
				originalDst.IP = net.IPv4(addr.Multiaddr[4], addr.Multiaddr[5], addr.Multiaddr[6], addr.Multiaddr[7])
			})
			if dstErr != nil {
				fmt.Printf("failed to obtain original dst: %v\n", dstErr)
				in.Close()
				return
			}
			out, err := net.DialTimeout("tcp", originalDst.String(), 5*time.Second)
			// out, err := net.Dial("tcp", originalDst.String())
			if err != nil {
				fmt.Printf("can not connect to original dst: %v, err: %v\n", originalDst.String(), err)
				return
			}
			defer out.Close()
			_ = out.SetDeadline(time.Now().Add(idleTimeout))
			go io.Copy(in, out)
			go io.Copy(out, in)
		}()
	}
}
