package middleware

import (
	"health_backend/pkg/metric"
	"time"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) gin.HandlerFunc {
	return func(g *gin.Context) {
		start := time.Now() // Start tracking time

		// Process the request
		g.Next()

		// Get status, method, and path
		status := g.Writer.Status()
		method := g.Request.Method
		path := g.FullPath() // Gets the registered route pattern

		// gapture response time
		duration := time.Since(start).Seconds()

		// Update Prometheus metrics
		metrics.ObserveResponseTime(status, method, path, duration)
		metrics.IncHits(status, method, path)
	}
}
