# gin-observability

项目整体基于gin框架，实现了对prometheus（metric）、otelgin（trace）、zap（log）的封装添加，达到开箱即用

### What it can do?

- metric：封装了三个基于gauge类型的指标，能够反馈api的请求次数、错误、请求耗时（R.E.D）
- trace：封装使用了otelgin官方的中间件，能够在gin.context传递trace信息
- log：封装了zap的初始化过程，自定义了输出格式，通过gin.context添加了traceID和spanID

### What it can't do?

- metric：仅仅作为客户端的exporter。不能进行数据的分析、告警与可视化，需配合prometheus和grafana等使用
- trace：不包含任何exporter。不能直接观测到链路信息，仅在日志中输出，若需要请参考otel文档并结合jaeger等使用
- log：只能结合gin.context使用。为了到达记录traceID和spanID的效果必须传入gin.context

### 使用

#### 集成使用

为了更加方便，引入了observability包将metric和trace模块整合到一起，可直接使用即可

``` go
import "github.com/Ltd5552/gin-observability/observability"

r := gin.New()
// 传入gin.engine和应用名（作为链路的服务名）
observability.Set(r, "ServerName")
defer log.Sync()
```

#### 分别使用

**若仅需要metric功能：**

``` go
import(
    "github.com/gin-gonic/gin"
    "github.com/Ltd5552/gin-observability/metric"
)

r := gin.New()

// 设置metric
metric.Set(r)
```

**若仅需要trace和log功能：**

``` go
import(
    "github.com/gin-gonic/gin"
    "github.com/Ltd5552/gin-observability/trace"
    "github.com/Ltd5552/gin-observability/log"
)

r := gin.New()

// 设置trace
trace.Set(r, "gin-observability")

defer log.Sync()

r.GET("/get", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "get",
    })
    // 使用log
    log.Info(c, "get successfully", zap.String("ping", "pong"))

})

```

**若三个都用到使用集成的observability模块**

``` go
import(
    "github.com/gin-gonic/gin"
    "github.com/Ltd5552/gin-observability/observability"
    "github.com/Ltd5552/gin-observability/log"
)

r := gin.New()

observability.Set(r, "gin-observability")

defer log.Sync()

r.GET("/get", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "get",
    })
    // 使用log
    log.Info(c, "get successfully", zap.String("ping", "pong"))

})

```

### P.S

三个模块均能单独使用，但trace和log相互依赖才有效果，而metric则完全独立

### What's the next?

设想：采用option模式，增添更多的可选择性：

- 是否添加trace exporter
- 是否自定义log日志文件参数