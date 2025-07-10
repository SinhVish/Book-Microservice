package main

import (
	"log"
	"net"
	"strconv"

	userGrpc "user-service/internal/grpc"
	"user-service/internal/repository"
	"user-service/internal/services"
	"user-service/pkg/config"
	"user-service/pkg/database"

	"shared/proto/user_service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	userRepo := repository.NewUserRepository(database.GetDB())
	userService := services.NewUserService(userRepo)
	userServer := userGrpc.NewUserServer(userService)

	log.Printf("Starting gRPC server in goroutine...")
	go startGRPCServer(userServer, cfg)

	log.Printf("Starting HTTP server...")
	startHTTPServer(cfg)
}

func startGRPCServer(userServer *userGrpc.UserServer, cfg *config.Config) {
	log.Printf("Starting gRPC server setup...")

	grpcPort, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Printf("Invalid port configuration: %v", err)
		return
	}
	grpcPort = grpcPort + 1000

	log.Printf("Attempting to listen on gRPC port %d", grpcPort)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Printf("Failed to listen on gRPC port: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	user_service.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("gRPC server starting on port %d", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Failed to start gRPC server: %v", err)
		return
	}
}

func startHTTPServer(cfg *config.Config) {
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

	log.Printf("HTTP server starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("Auth Service URL: %s", cfg.AuthServiceURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
