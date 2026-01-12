//go:build linux
// +build linux

package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func isConnClosedError(err error) bool {
	return errors.Is(err, syscall.ENOTCONN) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, syscall.ECONNRESET)
}

func transfer(dst *net.TCPConn, src *net.TCPConn, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(dst, src)
		if isConnClosedError(err) {
			return
		}
		dst.CloseWrite()
	}()
}

func getOriginalDst(conn *net.TCPConn) (*net.TCPAddr, error) {
	sc, err := conn.SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("getting syscall connection: %w", err)
	}
	originalDst := &net.TCPAddr{}
	var dstErr error
	sc.Control(func(fd uintptr) {
		addr, err := unix.GetsockoptIPv6Mreq(int(fd), unix.IPPROTO_IP, unix.SO_ORIGINAL_DST)
		if err != nil {
			dstErr = fmt.Errorf("getsockopt SO_ORIGINAL_DST: %w", err)
			return
		}
		originalDst.Port = int(addr.Multiaddr[2])<<8 + int(addr.Multiaddr[3])
		originalDst.IP = net.IPv4(addr.Multiaddr[4], addr.Multiaddr[5], addr.Multiaddr[6], addr.Multiaddr[7])
	})
	if dstErr != nil {
		return nil, fmt.Errorf("failed to obtain original dst: %w", dstErr)
	}
	return originalDst, nil
}

func getTCPConn(conn net.Conn) *net.TCPConn {
	tcpConn, _ := conn.(*net.TCPConn)
	tcpConn.SetDeadline(time.Now().Add(idleTimeout))
	return tcpConn
}
