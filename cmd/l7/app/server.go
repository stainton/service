//go:build linux
// +build linux

package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunHttpServer(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return fmt.Errorf("starting listener: %w", err)
	}
	engine := gin.Default()
	engine.Any("/*path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I got it!!!"})
	})
	return engine.RunListener(listener)
}
