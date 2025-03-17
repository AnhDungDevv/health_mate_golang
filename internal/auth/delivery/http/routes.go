package http

import (
	"health_backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handler) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/find", h.FindByName())
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:user_id", h.GetUserByID())
}
