//go:build linux
// +build linux

package app

import (
	"fmt"
	"io"
	"net"
	"time"
)

func ProxyHttp(input *IoWR, port int) error {
	output, err := net.DialTimeout("tcp", fmt.Sprintf(":%v", port), 5*time.Second)
	if err != nil {
		return fmt.Errorf("dialing to http server: %w", err)
	}
	defer output.Close()
	defer input.Close()

	go io.Copy(output, input)
	io.Copy(input, output)
	return nil
}
