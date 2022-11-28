package trace

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
)

func getIDMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if oteltrace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
			TraceID := oteltrace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
			SpanID := oteltrace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()
			fmt.Println(TraceID)
			fmt.Println(SpanID)
			logger.Debug("测试",
				zap.String("TraceID", TraceID),
				zap.String("SpanID", SpanID))
		}
	}
}

// InitTracer 初始化一个全局的TracerProvider
func initTracer() *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

// Set 封装了对InitTracer的调用，同时使用了两个中间件，r参数是gin的Engine，ServerName是自定义的服务名
func Set(r *gin.Engine, logger *zap.Logger, ServerName string) {
	// 获取一个全局TracerProvider
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	// 使用再封装的中间件
	r.Use(otelgin.Middleware(ServerName))
	r.Use(getIDMiddleware(logger))

}
