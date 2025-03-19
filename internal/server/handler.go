package server

import (
	"context"
	"health_backend/pkg/metric"
	"net/http"

	authHttp "health_backend/internal/auth/delivery/http"
	authRepository "health_backend/internal/auth/repository"
	authUseCase "health_backend/internal/auth/usecase"
	"health_backend/internal/chat/delivery/websocket"
	kafka_chat "health_backend/internal/chat/kafka"
	"health_backend/internal/chat/repository"
	chat_redis "health_backend/internal/chat/repository/redis"
	"health_backend/internal/chat/usecase"
	kafka "health_backend/internal/infrastructure/kafka"
	base_redis "health_backend/internal/infrastructure/redis"
	apiMiddlewares "health_backend/internal/middleware"

	"github.com/gin-gonic/gin"

	"health_backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) MapHandlers(g *gin.Engine) error {
	metrics, err := metric.NewMetrics(s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetric Error : %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	brokers := kafka.KafkaBrokers
	topicChat := kafka.TopicChat
	groupID := kafka.KafkaGroupID

	// Redis client is already initialized in Server struct
	baseRedis := base_redis.NewBaseRedisRepository(s.redisClient)
	chatRedisRepository := chat_redis.NewChatRedisRepository(baseRedis, s.logger)
	// init pgRepository
	chatRepository := repository.NewChatRepository(s.db)

	// Create context for Kafka
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize Kafka consumer and start it for chat module
	consumerChat := kafka_chat.NewKafkaConsumerChat(brokers, topicChat, groupID, s.logger, chatRedisRepository, chatRepository)
	// init chat usecae chat
	producerChat := kafka_chat.NewKafkaProducer(s.logger)
	chatUC := usecase.NewChatUsecase(s.cfg, chatRepository, chatRedisRepository, producerChat, s.logger)
	// init porducer

	// Start consumer in background with context
	go func() {
		consumerChat.StartConsumer()
		cancel()
	}()

	// Handle shutdown when context is done
	go func() {
		<-ctx.Done()
		s.logger.Info("Shutting down Kafka connections...")
		consumerChat.Close()
		producerChat.Close()
	}()

	// WebSocket Handler
	wsHandler := websocket.NewWebsocketHandler(s.cfg, chatUC, s.logger)

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHendler(s.cfg, authUC, s.logger)

	// Middlewares
	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, s.logger, metrics, authUC)

	// Swagger docs
	docs.SwaggerInfo.Title = "Go example REST API"
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// SSL redirect
	if s.cfg.Server.SSL {
		g.Use(func(c *gin.Context) {
			if c.Request.TLS == nil {
				c.Redirect(http.StatusMovedPermanently, "https://"+c.Request.Host+c.Request.RequestURI)
				c.Abort()
			}
		})
	}

	// CORS Middleware
	g.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-Request-ID, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	g.Use(
		mw.Recovery(),
		mw.RequestLoggerMiddleware(),
		mw.Metrics.Handler(),
	)
	// API Routes
	v1 := g.Group("/api/v1")
	// Auth routes
	authGroup := v1.Group("/auth")
	authHttp.MapAuthRoutes(authGroup, authHandlers)

	// WebSocket route - should be authenticated in production
	v1.GET("/ws", wsHandler.HandleWebSocket())

	return nil
}
