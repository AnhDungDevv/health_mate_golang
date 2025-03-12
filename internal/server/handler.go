package server

import (
	"fmt"
	"health_backend/pkg/metric"
	"net/http"

	"github.com/gin-gonic/gin"

	authHttp "health_backend/internal/auth/delivery/http"
	authRepository "health_backend/internal/auth/repository"
	authUseCase "health_backend/internal/auth/usecase"
	apiMiddlewares "health_backend/internal/middleware"

	"health_backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) MapHandlers(g *gin.Engine) error {
	metric, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	fmt.Println(metric)

	if err != nil {
		s.logger.Errorf("CreateMetric Error : %s", err)

	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)
	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHendler(s.cfg, authUC, s.logger)

	// Middelwares
	mv := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"}, s.logger)
	g.Use(mv.MetricsMiddleware(metric))
	// // Swagger docs
	docs.SwaggerInfo.Title = "Go example REST API"
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if s.cfg.Server.SSL {
		g.Use(func(g *gin.Context) {
			if g.Request.TLS != nil {
				g.Redirect(http.StatusMovedPermanently, "https://"+g.Request.Host+g.Request.RequestURI)
				g.Abort()
			}
		})
	}
	// CORS Middleware
	g.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	v1 := g.Group("/api/v1")
	authGroup := v1.Group("/auth")
	authHttp.MapAuthRoutes(authGroup, authHandlers)
	return nil

}
