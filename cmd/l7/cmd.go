//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"transparent/cmd/l7/app"
)

func main() {
	err := app.RunHttpServer(17001)
	if err != nil {
		fmt.Printf("server failed with error: %+v\n", err)
		os.Exit(1)
	}
}
