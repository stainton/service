package app

import (
	"bufio"
	"net"
)

type IoWR struct {
	*bufio.Reader
	*bufio.Writer
	handle net.Conn
}

func newWR(fromApp net.Conn) *IoWR {
	return &IoWR{
		Reader: bufio.NewReader(fromApp),
		Writer: bufio.NewWriter(fromApp),
		handle: fromApp,
	}
}

func (w *IoWR) Close() error {
	return w.handle.Close()
}
