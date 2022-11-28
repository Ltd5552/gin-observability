package main

import (
	"github.com/Ltd5552/gin-observability/log"
	"github.com/Ltd5552/gin-observability/observability"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	defer log.Sync()

	//// 设置metric
	//metric.Set(r)
	//
	//// 设置trace
	//trace.Set(r, "ServerName")

	observability.Set(r, "test")

	r.GET("/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get",
		})
		log.Info(c, "test")
	})

	_ = r.Run()
}
