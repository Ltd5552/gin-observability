package observability

import (
	"github.com/Ltd5552/gin-observability/log"
	"github.com/Ltd5552/gin-observability/metric"
	"github.com/Ltd5552/gin-observability/trace"
	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine, ServerName string) {
	// 设置metric
	metric.Set(r)

	// 设置trace
	trace.Set(r, ServerName)

	// 设置文件名
	log.FileName = ServerName + ".log"

}
