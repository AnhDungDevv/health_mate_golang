package main

// @title Go example REST API
// @version 1.0
// @description This is a sample server for a Go API with Swagger.
// @host localhost:8080
// @BasePath /api/v1
import (
	"health_backend/config"
	"health_backend/internal/server"
	"health_backend/pkg/db/postgres"
	"health_backend/pkg/db/redis"
	"health_backend/pkg/logger"
	"log"
)

func main() {
	log.Println("Starting API server...")

	// Load Configuration
	// configPath := utils.GetConfigPath(os.Getenv("CONFIG")) // Correct env var name
	cfgFile, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("LoadConfig error: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig error: %v", err)
	}

	// Initialize Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, Loglevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	// Initialize PostgreSQL Connection
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("Postgresql init failed: %s", err)
	} else {
		appLogger.Info("Postgres connected successfully")
	}

	// Get *sql.DB from *gorm.DB and defer closing
	sqlDB, err := psqlDB.DB()
	if err != nil {
		log.Fatalf("Error getting DB instance: %s", err)
	}
	defer sqlDB.Close()

	// Initialize Redis Client
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	s := server.NewServer(cfg, psqlDB, redisClient, appLogger)

	// Start Server
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
