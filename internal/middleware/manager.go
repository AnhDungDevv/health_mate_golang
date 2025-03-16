package middleware

import (
	"health_backend/config"
	"health_backend/internal/domain/usecase"
	"health_backend/internal/middleware"
	"health_backend/pkg/logger"
	"health_backend/pkg/metric"
)

type MiddlewareManager struct {
	cfg     *config.Config
	logger  logger.Logger
	metrics *middleware.MetricMiddleware
	authUC  usecase.AuthUseCase
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger, metrics metric.Metrics, authUC usecase.AuthUseCase) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:     cfg,
		logger:  logger,
		metrics: middleware.NewMetricMiddleware(metrics),
		authUC:  authUC,
	}
}
