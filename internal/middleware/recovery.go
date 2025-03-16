package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				mw.logger.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"path", c.Request.URL.Path,
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
