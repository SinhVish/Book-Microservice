package database

import (
	"fmt"
	"log"

	"user-service/internal/models"
	"user-service/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB is used to
// initialize the database connection and run migrations for user-service
// It's important because it establishes the database connection and creates necessary tables
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

// AutoMigrate is used to
// automatically create or update database tables based on model definitions
// It's important because it ensures the database schema matches the Go models
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
	)
}

// GetDB is used to
// get the database instance for use in other parts of the application
// It's important because it provides access to the database connection for repositories
func GetDB() *gorm.DB {
	return DB
}
