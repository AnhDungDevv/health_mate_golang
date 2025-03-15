package server

import (
	"health_backend/config"
	"health_backend/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	gin         *gin.Engine
	cfg         *config.Config
	db          *gorm.DB
	redisClient *redis.Client
	logger      logger.Logger
}

func NewServer(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, logger logger.Logger) *Server {
	return &Server{
		gin:         gin.Default(),
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (s *Server) Run() error {
	err := s.MapHandlers(s.gin)
	if err != nil {
		s.logger.Errorf("Lỗi trong MapHandlers: %v", err)
		return err
	}

	server := &http.Server{
		Addr:         "0.0.0.0" + s.cfg.Server.Port, // ĐÚNG
		Handler:      s.gin,
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error starting Server: %v", err)
		}
	}()

	// go func() {
	// 	s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
	// 	if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
	// 		s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
	// 	}
	// }()

	return nil
}
