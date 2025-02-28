package http

import (
	"health_backend/config"
	auth "health_backend/internal/auth/interfaces"
	"health_backend/internal/models"
	"health_backend/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

// Delete implements auth.Handler.
func (a *authHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// FindByName implements auth.Handler.
func (a *authHandlers) FindByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// GetCSRFToken implements auth.Handler.
func (a *authHandlers) GetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// GetMe implements auth.Handler.
func (a *authHandlers) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// GetUserByID implements auth.Handler.
func (a *authHandlers) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// GetUsers implements auth.Handler.
func (a *authHandlers) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// Login implements auth.Handler.
func (a *authHandlers) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// Logout implements auth.Handler.
func (a *authHandlers) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Function not implement"})
	}
}

// Register implements auth.Handler.
// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
func (a *authHandlers) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		user := &models.User{}
		if err := c.ShouldBindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdUser, err := a.authUC.Register(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot register user"})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

// Update implements auth.Handler.
func (a *authHandlers) Update() gin.HandlerFunc {
	panic("unimplemented")
}

// UploadAvatar implements auth.Handler.
func (a *authHandlers) UploadAvatar() gin.HandlerFunc {
	panic("unimplemented")
}

func NewAuthHendler(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handler {
	return &authHandlers{
		cfg: cfg, authUC: authUC, logger: log,
	}
}
