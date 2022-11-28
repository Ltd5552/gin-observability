package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

// 创建Metric数据结构
var (
	RequestDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "RequestDuration",
		Help: "请求耗时(ms)",
	},
		// 添加标签为请求路径和方法
		[]string{"path", "method"})

	RequestCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "RequestCount",
		Help: "请求次数",
	},
		// 添加标签为请求路径、方法和状态码
		[]string{"path", "method", "code"})

	RequestErrorCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "RequestErrorCount",
		Help: "请求错误次数",
	},
		// 添加标签为请求路径和方法
		[]string{"path", "method"})
)

// 初始化
func init() {
	// register the three gauges
	prometheus.MustRegister(RequestDuration, RequestCount, RequestErrorCount)
}

func PrometheusMW(c *gin.Context) {
	// 开始计时
	Start := time.Now()

	// after request
	c.Next()

	// 计算请求的耗时
	Duration := float64(time.Since(Start) / time.Millisecond)

	// 获取请求的路径,请求方式,状态码
	Path := c.FullPath()
	Method := c.Request.Method
	StatusCode := strconv.Itoa(c.Writer.Status())

	Error := c.Errors
	if Error != nil {
		// 向RequestErrorCount中写入
		RequestErrorCount.WithLabelValues(Path, Method).Inc()
	}
	// 向RequestCount中写入
	RequestCount.WithLabelValues(Path, Method, StatusCode).Inc()
	// 向RequestDuration中写入
	RequestDuration.WithLabelValues(Path, Method).Set(Duration)
}

func Set(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Use(PrometheusMW)
}
