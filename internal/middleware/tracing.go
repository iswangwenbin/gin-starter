package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TraceIDKey = "trace_id"
	SpanIDKey  = "span_id"
)

// TracingMiddleware 链路追踪中间件
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成或获取 trace ID
		traceID := c.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		
		// 生成 span ID
		spanID := uuid.New().String()
		
		// 设置到 context
		c.Set(TraceIDKey, traceID)
		c.Set(SpanIDKey, spanID)
		
		// 设置响应头
		c.Header("X-Trace-Id", traceID)
		c.Header("X-Span-Id", spanID)
		
		// 记录请求开始
		start := time.Now()
		
		// 处理请求
		c.Next()
		
		// 记录请求结束
		duration := time.Since(start)
		
		// 记录访问日志
		logger := zap.L().With(
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
		
		if c.Writer.Status() >= 400 {
			logger.Error("Request completed with error")
		} else {
			logger.Info("Request completed successfully")
		}
	}
}

// GetTraceID 从 context 获取 trace ID
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		return traceID.(string)
	}
	return ""
}

// GetSpanID 从 context 获取 span ID
func GetSpanID(c *gin.Context) string {
	if spanID, exists := c.Get(SpanIDKey); exists {
		return spanID.(string)
	}
	return ""
}

// WithTraceContext 将 trace 信息添加到 context
func WithTraceContext(c *gin.Context, ctx context.Context) context.Context {
	if traceID := GetTraceID(c); traceID != "" {
		ctx = context.WithValue(ctx, TraceIDKey, traceID)
	}
	if spanID := GetSpanID(c); spanID != "" {
		ctx = context.WithValue(ctx, SpanIDKey, spanID)
	}
	return ctx
}