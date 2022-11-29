package observability

import (
	"github.com/Ltd5552/gin-observability/metric"
	"github.com/Ltd5552/gin-observability/trace"
	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine, ServerName string) {
	// 设置metric
	metric.Set(r)

	// 设置trace
	trace.Set(r, ServerName)

	// log的自定义？
}
