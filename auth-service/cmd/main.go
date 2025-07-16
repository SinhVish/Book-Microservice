package main

import (
	"log"
	"net"
	"strconv"

	"auth-service/internal/clients"
	authGrpc "auth-service/internal/grpc"
	"auth-service/internal/handlers"
	"auth-service/internal/repository"
	"auth-service/internal/services"
	"auth-service/pkg/config"
	"auth-service/pkg/database"

	"shared/proto/auth_service"

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

	userServiceClient, err := clients.NewUserServiceClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatal("Failed to create user service client:", err)
	}
	defer userServiceClient.Close()

	credentialRepo := repository.NewCredentialRepository(database.GetDB())

	authService := services.NewAuthService(credentialRepo, userServiceClient, cfg)

	authHandler := handlers.NewAuthHandler(authService)

	authServer := authGrpc.NewAuthServer(cfg)
	go startGRPCServer(authServer, cfg)

	startHTTPServer(cfg, authHandler)
}

func startGRPCServer(authServer *authGrpc.AuthServer, cfg *config.Config) {
	log.Printf("Starting gRPC server setup...")

	grpcPort, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Printf("Invalid port configuration: %v", err)
		return
	}
	grpcPort = grpcPort + 1000 // Use port 9080 for gRPC if HTTP is 8080

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Printf("Failed to listen on gRPC port: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	auth_service.RegisterAuthServiceServer(grpcServer, authServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Failed to start gRPC server: %v", err)
		return
	}
}

func startHTTPServer(cfg *config.Config, authHandler *handlers.AuthHandler) {
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

	log.Printf("HTTP server starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("User Service URL: %s", cfg.UserServiceURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
