# Log

## 介绍

基于zap的封装，为了在日志中添加traceID和spanID两个信息每次记录日志时需传递gin.context进去，需配合trace包使用

## 使用
由于存在是否需要传入gin.context两种情况，所以保留了初始的方法Info，而在需要传入gin.context的方法后面统一加了`WithID`，比如InfoWithID等等

``` go
import (
	"github.com/Ltd5552/gin-observability/log"
	"github.com/Ltd5552/gin-observability/trace"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	defer log.Sync()

	// 添加trace
	trace.Set(r, "test")

	r.GET("/get", func(c *gin.Context) {
        // 使用log，第一个参数是gin.Context
		log.InfoWithID(c, "test")
		log.Info("test")
	})
}

```

### 结果

![log](../pic/log.png)