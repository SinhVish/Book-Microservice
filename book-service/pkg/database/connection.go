package database

import (
	"log"

	"book-service/internal/models"
	"book-service/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error

	dsn := cfg.GetDatabaseURL()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Println("Connected to database successfully")

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.Author{}, &models.Book{}); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")

	return nil
}

func GetDB() *gorm.DB {
	return db
}
