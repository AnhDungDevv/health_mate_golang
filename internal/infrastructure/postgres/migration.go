package postgres_migration

import (
	"health_backend/internal/models"
	"log"

	"gorm.io/gorm"
)

func MigrateAll(db *gorm.DB) {
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Expertiese{},
	}
	err := db.AutoMigrate(modelsToMigrate...)
	if err != nil {
		log.Fatalf("Migration failed: %s", err)
	}

	log.Println("✅ Database migration completed successfully.")
}
