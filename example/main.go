package main

import (
	"github.com/Ltd5552/gin-observability/log"
	"github.com/Ltd5552/gin-observability/metric"
	"github.com/Ltd5552/gin-observability/trace"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	defer log.Sync()

	// 测试metric
	metric.Set(r)

	// 测试trace
	trace.Set(r, "test")

	r.GET("/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get",
		})
		log.Info(c, "test")
	})

	_ = r.Run()
}
