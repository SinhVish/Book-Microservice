package database

import (
	"fmt"
	"log"

	"auth-service/internal/models"
	"auth-service/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := cfg.GetDatabaseURL()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	DB = db

	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Credential{},
		&models.RefreshToken{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
