# Metric

## 介绍

基于prometheus的client_golang和gin实现，该中间件实现了对gin框架下api的三个指标，包含请求耗时、请求次数以及请求错误，基于promQL能实现R.E.D监控指标，默认/metrics路径为prometheus获取


## 使用

```
import "github.com/Ltd5552/gin-observability/metric"
```

```
// 初始化gin
r := gin.Default()

// 将gin.Engine——r传入Set()即可
metric.Set(r)

// your code ...
```

### 快速开始

下面是一个示例

``` go
package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/Ltd5552/gin-observability/metric"
)

var (
	TestError = fmt.Errorf("just a test")
)

func main() {
	r := gin.Default()

	metric.Set(r)
    
	r.GET("/get", func(context *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		context.JSON(200, gin.H{
			"message": "get",
		})
	})

	r.POST("/post", func(context *gin.Context) {
		time.Sleep(66 * time.Millisecond)
		context.JSON(200, gin.H{
			"message": "post",
		})
	})

	r.DELETE("/delete", func(context *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		context.JSON(200, gin.H{
			"message": "delete",
		})
	})

	r.PUT("/put", func(context *gin.Context) {
		time.Sleep(50 * time.Millisecond)
		context.JSON(200, gin.H{
			"message": "put",
		})
	})

	r.GET("/error", func(context *gin.Context) {
		context.Error(TestError)
	})

	err := r.Run(":80")  // your port 
	if err != nil {
		return
	}
}
```



## 显示

`http://localhost/metrics`

里面默认包含了一些基本的go相关信息，而api相关信息会在该api第一次请求后显示
