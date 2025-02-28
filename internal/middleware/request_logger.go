package middleware

import (
	"health_backend/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) RequestLoggerMiddleware(next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		// get request and response from context
		req := ctx.Request
		res := ctx.Writer

		status := res.Status()
		size := res.Size()
		elapsedTime := time.Since(start).String()
		requestID := utils.GetRequestID(ctx)

		mw.logger.Infof("RequestID: %s, Method: %s, URI: %s, Status: %d, Size: %d, Time: %s",
			requestID, req.Method, req.URL.Path, status, size, elapsedTime,
		)
	}
}
