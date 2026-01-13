package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"syscall"

	"github.com/gin-gonic/gin"
)

type fakeListener struct {
	conn        *net.TCPConn
	isConnected bool
}

func isHTTP(conn io.Reader) bool {
	reader := bufio.NewReader(conn)
	head, err := reader.Peek(64)
	if err != nil {
		return false
	}
	return bytes.HasPrefix(head, []byte("GET ")) ||
		bytes.HasPrefix(head, []byte("POST ")) ||
		bytes.HasPrefix(head, []byte("HEAD ")) ||
		bytes.HasPrefix(head, []byte("PUT ")) ||
		bytes.HasPrefix(head, []byte("DELETE ")) ||
		bytes.HasPrefix(head, []byte("OPTIONS ")) ||
		bytes.HasPrefix(head, []byte("PATCH ")) ||
		bytes.HasPrefix(head, []byte("CONNECT ")) ||
		bytes.HasPrefix(head, []byte("TRACE ")) ||
		bytes.HasPrefix(head, []byte("HTTP/1."))
}

func (fl *fakeListener) Accept() (net.Conn, error) {
	if fl.isConnected {
		return nil, syscall.EINVAL
	}
	fl.isConnected = true
	return fl.conn, nil
}

func (fl *fakeListener) Close() error {
	fl.isConnected = false
	return fl.conn.Close()
}

func (fl *fakeListener) Addr() net.Addr {
	return fl.conn.LocalAddr()
}

func toFakeListener(conn *net.TCPConn, isConnected bool) net.Listener {
	return &fakeListener{conn: conn, isConnected: isConnected}
}

func serveAsHTTP(conn *net.TCPConn, dst string) {

	// 创建反向代理
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// 修改请求指向真实目标
			req.URL.Scheme = "http"
			req.URL.Host = dst
			req.RequestURI = ""
			fmt.Printf("[HTTP] Forward %s %s -> %s\n", req.Method, req.RequestURI, dst)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			fmt.Printf("[HTTP] proxy error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Bad Gateway"))
		},
	}

	engine := gin.Default()
	engine.Any("/*path", func(ctx *gin.Context) {
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})

	engine.RunListener(toFakeListener(conn, false))

}
