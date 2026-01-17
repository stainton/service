//go:build linux
// +build linux

package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

func ProxyTcp(port, httpPort int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return fmt.Errorf("listening on port %d: %w", port, err)
	}
	defer listener.Close()
	fmt.Printf("listener start listening\n")
	// Accept connections and handle them
	for {
		newConn, err := listener.Accept()
		fmt.Printf("accept new connection\n")
		if err != nil {
			return fmt.Errorf("accepting connection: %w", err)
		}

		fromApp := newWR(newConn)

		if isHttp(fromApp.Reader) {
			go ProxyHttp(fromApp, httpPort)
			continue
		}
		go func() {
			defer fromApp.Close()
			addr, err := getOriginalDst(newConn)
			if err != nil {
				fmt.Printf("getting original destination: %v\n", err)
				return
			}
			dst, err := net.Dial("tcp", addr.String())
			if err != nil {
				fmt.Printf("dialing to original destination %s: %v\n", addr.String(), err)
				return
			}
			defer dst.Close()
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				_, err := io.Copy(dst, fromApp)
				if err != nil {
					fmt.Printf("copying from app to dst: %v\n", err)
				}
			}()
			go func() {
				defer wg.Done()
				_, err := io.Copy(fromApp, dst)
				if err != nil {
					fmt.Printf("copying from dst to app: %v\n", err)
				}
			}()
			wg.Wait()
		}()
	}
}

func isHttp(conn *bufio.Reader) bool {
	peek, err := conn.Peek(64)
	if err != nil {
		return false
	}
	return bytes.HasPrefix(peek, []byte("GET ")) ||
		bytes.HasPrefix(peek, []byte("POST ")) ||
		bytes.HasPrefix(peek, []byte("HEAD ")) ||
		bytes.HasPrefix(peek, []byte("PUT ")) ||
		bytes.HasPrefix(peek, []byte("DELETE ")) ||
		bytes.HasPrefix(peek, []byte("OPTIONS ")) ||
		bytes.HasPrefix(peek, []byte("CONNECT ")) ||
		bytes.HasPrefix(peek, []byte("TRACE ")) ||
		bytes.HasPrefix(peek, []byte("HTTP/1. "))
}

func getOriginalDst(conn net.Conn) (*net.TCPAddr, error) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("not a TCP connection")
	}
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("getting syscall.RawConn: %w", err)
	}
	originalDst := &net.TCPAddr{}
	var sockerr error
	rawConn.Control(func(fd uintptr) {
		addr, err := syscall.GetsockoptIPv6Mreq(int(fd), unix.IPPROTO_IP, unix.SO_ORIGINAL_DST)
		if err != nil {
			sockerr = err
			return
		}
		ip := net.IPv4(addr.Multiaddr[4], addr.Multiaddr[5], addr.Multiaddr[6], addr.Multiaddr[7])
		port := int(addr.Multiaddr[2])<<8 + int(addr.Multiaddr[3])
		originalDst = &net.TCPAddr{IP: ip, Port: port}
	})
	return originalDst, sockerr
}
