package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port    string
	GinMode string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	AuthServiceURL string
	RedisURL       string

	// Cache configuration
	CacheEnabled      bool
	L1CacheTTLMinutes int
	L2CacheTTLMinutes int

	JWTSecret string
}

func LoadConfig() (*Config, error) {
	// Parse cache configuration
	cacheEnabled, _ := strconv.ParseBool(getEnv("CACHE_ENABLED", "true"))
	l1CacheTTL, _ := strconv.Atoi(getEnv("L1_CACHE_TTL_MINUTES", "5"))
	l2CacheTTL, _ := strconv.Atoi(getEnv("L2_CACHE_TTL_MINUTES", "15"))

	config := &Config{
		Port:    getEnv("PORT", "8081"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "user_db"),

		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "localhost:9080"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),

		CacheEnabled:      cacheEnabled,
		L1CacheTTLMinutes: l1CacheTTL,
		L2CacheTTLMinutes: l2CacheTTL,

		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
	}

	return config, nil
}

func (c *Config) GetL1CacheTTL() time.Duration {
	return time.Duration(c.L1CacheTTLMinutes) * time.Minute
}

func (c *Config) GetL2CacheTTL() time.Duration {
	return time.Duration(c.L2CacheTTLMinutes) * time.Minute
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
