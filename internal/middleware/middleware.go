package middleware

import (
	"health_backend/config"
	auth "health_backend/internal/auth/interfaces"
	"health_backend/pkg/logger"
)

type MiddlewareManager struct {
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

func NewMiddlewareManager(authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, origins: origins, logger: logger}

}
