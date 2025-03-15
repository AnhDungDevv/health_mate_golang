package server

import (
	"fmt"
	"health_backend/pkg/metric"
	"log"
	"net/http"

	authHttp "health_backend/internal/auth/delivery/http"
	authRepository "health_backend/internal/auth/repository"
	authUseCase "health_backend/internal/auth/usecase"
	apiMiddlewares "health_backend/internal/middleware"
	"health_backend/internal/mqtt/delivery"
	mqttRoutuer "health_backend/internal/mqtt/delivery"
	repositoryMQTT "health_backend/internal/mqtt/repository"
	"health_backend/internal/mqtt/usecase"

	"github.com/gin-gonic/gin"

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
	// Init MQTT Client
	mqttClient := repositoryMQTT.NewMQTTClient("tcp://emqx:1883", "your-client-id")
	err = mqttClient.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	mqttUC := usecase.NewMQTTUsecase(mqttClient)

	// Init handlers
	authHandlers := authHttp.NewAuthHendler(s.cfg, authUC, s.logger)
	mqttHandler := delivery.NewMQTTHandler(mqttUC)

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
	mqttRoutuer.NewMapMQTTRoutes(v1, mqttHandler)
	authHttp.MapAuthRoutes(authGroup, authHandlers)
	return nil

}
