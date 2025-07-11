package main

import (
	"log"

	"auth-service/internal/clients"
	"auth-service/internal/handlers"
	"auth-service/internal/repository"
	"auth-service/internal/services"
	"auth-service/pkg/config"
	"auth-service/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	userServiceClient, err := clients.NewUserServiceClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatal("Failed to create user service client:", err)
	}
	defer userServiceClient.Close()

	credentialRepo := repository.NewCredentialRepository(database.GetDB())
	authService := services.NewAuthService(credentialRepo, userServiceClient, cfg)
	authHandler := handlers.NewAuthHandler(authService)

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
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	log.Printf("Auth service starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("User Service URL: %s", cfg.UserServiceURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
