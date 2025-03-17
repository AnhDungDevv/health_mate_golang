package http

import (
	"health_backend/config"
	"health_backend/internal/auth"
	"health_backend/internal/models"
	httpErrors "health_backend/pkg/httpError"
	"health_backend/pkg/logger"
	"health_backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type Handler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	FindByName() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	GetMe() gin.HandlerFunc
	UploadAvatar() gin.HandlerFunc
}

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

// Register godoc
// @Summary Register new user
// @Description Create a new user account and return user details with authentication tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} RegisterResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 409 {object} ErrorResponse "Email already in use"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// Login godoc
// @Summary user login
// @Description Authenticate user with email/username and password, returning access & refresh tokens
// @Tags Authentication
// @Accept  json
// @Produce json
// @Params request body LoginRequest true "Login credentials"
// @Success  200 {object} LoginResponse "Successful login"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *authHandlers) Login() gin.HandlerFunc {
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password`
	}
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.login")
		defer span.Finish()

		login := &Login{}

		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			statusCode, responseBody := httpErrors.ErrorResponse(err)
			c.JSON(statusCode, responseBody)
		}

		userWithToken, err := h.authUC.Login(ctx, &models.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			statusCode, responseBody := httpErrors.ErrorResponse(err)
			c.JSON(statusCode, responseBody)
		}
		c.JSON(http.StatusOK, userWithToken)

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
