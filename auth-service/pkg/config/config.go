package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port    string
	GinMode string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret      string
	JWTExpiryHours int

	UserServiceURL string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "auth_db"),

		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),

		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
	}

	jwtExpiryStr := getEnv("JWT_EXPIRY_HOURS", "24")
	jwtExpiry, err := strconv.Atoi(jwtExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %v", err)
	}
	config.JWTExpiryHours = jwtExpiry

	return config, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
