package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) AuthMiddlewre() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		// Remove Bearer prefix
		token = strings.Replace(token, "Bearer ", "", 1)

		// Validate token
		userID, err := mw.authUC.ValidateToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// Set user info to context
		c.Set("user_id", userID)
		c.Next()
	}
}
