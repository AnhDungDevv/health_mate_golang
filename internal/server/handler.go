package server

import (
	"net/http"

	// Import Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"health_backend/docs"
	authHttp "health_backend/internal/auth/delivery/http"
	authRepository "health_backend/internal/auth/repository"
	authUseCase "health_backend/internal/auth/usecase"
	"health_backend/pkg/metric"

	"github.com/gin-gonic/gin"
	// "github.com/swaggo/swag/example/override/docs"
)

func (s *Server) MapHandlers(g *gin.Engine) error {
	metric, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)

	if err != nil {
		s.logger.Errorf("CreateMetric Error : %s", err)

	}
	// s.logger.Info(
	// 	"Metrics available URL: %s, ServiceName: %s",
	// 	s.cfg.Metrics.URL,
	// 	s.cfg.Metrics.ServiceName,
	// )
	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHendler(s.cfg, authUC, s.logger)

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
		c.Header("Access-Control-Allow-Origin", "*") // Cho phép tất cả hoặc thay bằng domain cụ thể
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // Trả về 204 thay vì 200
			return
		}

		c.Next()
	})

	v1 := g.Group("/api/v1")
	authGroup := v1.Group("/auth")
	authHttp.MapAuthRoutes(authGroup, authHandlers)
	return nil

}
