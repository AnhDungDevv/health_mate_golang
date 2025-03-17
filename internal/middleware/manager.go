package middleware

import (
	"health_backend/config"
	"health_backend/internal/auth"
	"health_backend/pkg/logger"
	"health_backend/pkg/metric"
)

type MiddlewareManager struct {
	cfg     *config.Config
	logger  logger.Logger
	Metrics *MetricMiddleware
	authUC  auth.UseCase
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger, metrics metric.Metrics, authUC auth.UseCase) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:     cfg,
		logger:  logger,
		Metrics: NewMetricMiddleware(metrics),
		authUC:  authUC,
	}
}
