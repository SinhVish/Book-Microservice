package main

import (
	"log"

	"user-service/pkg/config"

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
			"service": "user-service",
			"port":    cfg.Port,
		})
	})

	userGroup := r.Group("/api/v1/users")
	{
		userGroup.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "User profile endpoint - coming soon",
			})
		})
		userGroup.PUT("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Update user profile endpoint - coming soon",
			})
		})
		userGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "List users endpoint - coming soon",
			})
		})
	}

	log.Printf("User service starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("Auth Service URL: %s", cfg.AuthServiceURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
