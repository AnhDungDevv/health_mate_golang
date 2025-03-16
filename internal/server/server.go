package server

import (
	"context"
	"health_backend/config"
	"health_backend/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	httpServer  *http.Server
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
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	err := s.MapHandlers(s.gin)
	if err != nil {
		s.logger.Errorf("Lỗi trong MapHandlers: %v", err)
		return err
	}

	s.httpServer = &http.Server{
		Addr:         "0.0.0.0" + s.cfg.Server.Port,
		Handler:      s.gin,
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error starting Server: %v", err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	s.logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Errorf("Server forced to shutdown: %v", err)
	}

	s.logger.Info("Server exiting")
	return nil
}

// Cleanup performs cleanup operations
func (s *Server) Cleanup() {
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}
