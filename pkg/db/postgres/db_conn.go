package postgres

import (
	"fmt"
	"health_backend/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Các thông số kết nối
const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// NewPsqlDB - Connect PostgreSQL by GORM
func NewPsqlDB(p *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		p.Postgres.PostgresqlHost,
		p.Postgres.PostgresqlPort,
		p.Postgres.PostgresqlUser,
		p.Postgres.PostgresqlDbname,
		p.Postgres.PostgresqlPassword,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //  log SQL
	})
	if err != nil {
		return nil, err
	}

	// Get instance SQL to config for connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Second)

	// Check the connection
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Successfully connected to PostgreSQL with GORM")
	return db, nil
}
