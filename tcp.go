package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

const idleTimeout = 10 * time.Second

func TCPServer() {
	listener, _ := net.Listen("tcp", ":15001")
	for {
		in, _ := listener.Accept()
		go func() {
			fmt.Println("[Debug] new connection coming")
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
				originalDst.Port = int(addr.Multiaddr[2])<<8 + int(addr.Multiaddr[3])
				originalDst.IP = net.IPv4(addr.Multiaddr[4], addr.Multiaddr[5], addr.Multiaddr[6], addr.Multiaddr[7])
			})
			if dstErr != nil {
				fmt.Printf("failed to obtain original dst: %v\n", dstErr)
				in.Close()
				return
			}
			fmt.Printf("[Debug] Dialing %v:%v\n", originalDst.IP.String(), originalDst.Port)
			out, err := net.DialTimeout("tcp", originalDst.String(), 5*time.Second)
			if err != nil {
				fmt.Printf("can not connect to original dst: %v, err: %v\n", originalDst.String(), err)
				return
			}
			defer out.Close()
			_ = in.SetDeadline(time.Now().Add(idleTimeout))
			_ = out.SetDeadline(time.Now().Add(idleTimeout))
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				io.Copy(in, out)
				err := in.(*net.TCPConn).CloseWrite()
				if err != nil {
					fmt.Printf("close write to in error: %v\n", err)
				}
				fmt.Printf("finish transfer from out to in\n")
			}()
			go func() {
				defer wg.Done()
				io.Copy(out, in)
				err := out.(*net.TCPConn).CloseWrite()
				if err != nil {
					fmt.Printf("close write to out error: %v\n", err)
				}
				fmt.Printf("finish transfer from in to out\n")
			}()
			wg.Wait()
			fmt.Printf("connection from %v to %v closed\n", in.RemoteAddr().String(), originalDst.String())
		}()
	}
}
