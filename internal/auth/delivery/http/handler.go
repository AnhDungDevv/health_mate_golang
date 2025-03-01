package http

import (
	"health_backend/config"
	auth "health_backend/internal/auth/interfaces"
	"health_backend/internal/models"
	httpErrors "health_backend/pkg/httpError"
	"health_backend/pkg/logger"
	"health_backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHendler(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handler {
	return &authHandlers{
		cfg: cfg, authUC: authUC, logger: log,
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
func (h *authHandlers) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Register")
		defer span.Finish()

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			c.JSON(httpErrors.ErrorResponse(err))
			return
		}

		createdUser, err := h.authUC.Register(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			c.JSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(http.StatusCreated, createdUser)

	}
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

// Update implements auth.Handler.
func (a *authHandlers) Update() gin.HandlerFunc {
	panic("unimplemented")
}

// UploadAvatar implements auth.Handler.
func (a *authHandlers) UploadAvatar() gin.HandlerFunc {
	panic("unimplemented")
}
