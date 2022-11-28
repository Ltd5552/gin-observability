# gin-observability

项目整体基于gin框架，实现了对prometheus（metric）、otelgin（trace）、zap（log）的封装添加，达到开箱即用

### What it can do?

- metric：封装了三个基于gauge类型的指标，能够反馈api的请求次数、错误、请求耗时（R.E.D）
- trace：封装使用了otelgin官方的中间件，能够在gin.context传递trace信息
- log：封装了zap的初始化过程，自定义了输出格式，通过gin.context添加了traceID和spanID

### What it can't do?

- metric：仅仅作为客户端的exporter。不能进行数据的分析、告警与可视化，需配合promethues和grafana等使用
- trace：不包含任何exporter。不能直接观测到链路信息，仅在日志中输出，若需要请参考otel文档并结合jaeger等使用
- log：只能结合gin.context使用。为了到达记录traceaID和spanID的效果必须传入gin.context



### 使用

``` go
import "github.com/Ltd5552/gin-observability/observability"

// 三个模块集成
observability.Set(r, "test")
```



### P.S

三个模块均能单独使用，但trace和log相互依赖，而metric则完全独立