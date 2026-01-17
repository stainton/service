//go:build linux
// +build linux

package l4

import (
	"fmt"
	"os"
	"transparent/cmd/l4/app"
)

func main() {
	err := app.ProxyTcp(15001, 17001)
	if err != nil {
		fmt.Printf("proxy failed with error: %+v\n", err)
		os.Exit(1)
	}
}
