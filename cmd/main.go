//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	l4 "transparent/cmd/l4/app"
	l7 "transparent/cmd/l7/app"
)

func main() {
	tcpProxyPort := 15001
	httpListenPort := 17001
	go func() {
		err := l7.RunHttpServer(httpListenPort)
		if err != nil {
			panic(err)
		}
	}()
	err := l4.ProxyTcp(tcpProxyPort, httpListenPort)
	if err != nil {
		fmt.Printf("proxy failed with error: %+v\n", err)
		os.Exit(1)
	}
}
