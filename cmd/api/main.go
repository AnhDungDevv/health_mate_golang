package main

import (
	"health_backend/config"
	"health_backend/internal/server"
	db "health_backend/migration"
	"health_backend/pkg/db/postgres"
	"health_backend/pkg/db/redis"
	"health_backend/pkg/logger"
	"log"
)

// @title Health Backend API
// @version 1.0
// @description API cho hệ thống Health Backend
// @host localhost:5000
// @BasePath /api/v1
func main() {
	log.Println("Starting API server...")

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
	db.MigrateAll(psqlDB)
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
