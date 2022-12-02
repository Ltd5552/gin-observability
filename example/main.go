package main

import (
	"github.com/Ltd5552/gin-observability/log"
	"github.com/Ltd5552/gin-observability/observability"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	r := gin.New()

	// 设置metric
	//metric.Set(r)

	// 设置trace
	//trace.Set(r, "ServerName")

	observability.Set(r, "gin-observability")

	defer log.Sync()

	r.GET("/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get",
		})
		log.Info(c, "get successfully", zap.String("ping", "pong"))

	})

	_ = r.Run()
}
