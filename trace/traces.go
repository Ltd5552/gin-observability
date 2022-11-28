package trace

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

// InitTracer 初始化一个全局的TracerProvider，不含exporter，只将在log中打印出来，未添加可观测性后端
func initTracer() *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

// Set 封装了对InitTracer的调用，同时使用了两个中间件，r参数是gin的Engine，ServerName是自定义的服务名
func Set(r *gin.Engine, ServerName string) {
	// 获取一个全局TracerProvider
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	// 使用open-telemetry官方中间件
	r.Use(otelgin.Middleware(ServerName))
}
