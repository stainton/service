package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpServer() {
	engine := gin.Default()

	engine.Any("/*path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "你的请求，我劫持了！"})
	})
	engine.Run(":16001")
}
