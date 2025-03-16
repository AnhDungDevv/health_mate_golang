package middleware

import (
	"health_backend/pkg/metric"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricMiddleware struct {
	metrics metric.Metrics
}

func NewMetricMiddleware(metrics metric.Metrics) *MetricMiddleware {
	return &MetricMiddleware{
		metrics: metrics,
	}
}

func (m *MetricMiddleware) MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		m.metrics.IncHits(status, method, path)
		m.metrics.ObserveResponseTime(status, method, path, duration)
	}
}
