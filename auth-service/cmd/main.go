package main

import (
	"log"

	"auth-service/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "auth-service",
			"port":    cfg.Port,
		})
	})

	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Login endpoint - coming soon",
			})
		})
		authGroup.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Register endpoint - coming soon",
			})
		})
		authGroup.POST("/validate", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Token validation endpoint - coming soon",
			})
		})
	}

	log.Printf("Auth service starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("User Service URL: %s", cfg.UserServiceURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
