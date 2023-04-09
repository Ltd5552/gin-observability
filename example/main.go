package main

import (
	"github.com/Ltd5552/gin-observability/log"
	"github.com/Ltd5552/gin-observability/observability"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	// 直接统一设置
	observability.Set(r, "gin-observability")

	defer log.Sync()

	r.GET("/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get",
		})
		// 使用普通的zap.Info
		log.Info("get successfully")

		// 使用InfoWithID，第一个参数是gin.Context
		log.InfoWithID(c, "get successfully")
	})

	_ = r.Run()
}

//Happy B-day Goland
