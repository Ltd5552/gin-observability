# Trace
## 介绍

基于open-telemetry-gin的封装，没有使用exporter，需结合[log](https://github.com/Ltd5552/gin-observability/tree/main/log)包使用，信息记录在context后续日志中打印

若不使用log将没有任何反馈，但你依然可以在代码中通过gin.Context中主动获取到

## 使用

``` go
import "github.com/Ltd5552/gin-observability/trace"
```

``` go
// 初始化gin	
r := gin.New()
// 传入gin.Engine和项目名
trace.Set(r, "test")
// your code...
```



### 主动获取

``` go
import oteltrace "go.opentelemetry.io/otel/trace"

if oteltrace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
		traceID = oteltrace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
		spanID = oteltrace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()
	}
```

