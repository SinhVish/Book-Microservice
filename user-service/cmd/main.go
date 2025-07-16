package main

import (
	"log"
	"net"
	"strconv"

	"user-service/internal/clients"
	userGrpc "user-service/internal/grpc"
	"user-service/internal/handlers"
	"user-service/internal/middleware"
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

	authServiceClient, err := clients.NewAuthServiceClient(cfg.AuthServiceURL)
	if err != nil {
		log.Fatal("Failed to create auth service client:", err)
	}
	defer authServiceClient.Close()

	userRepo := repository.NewUserRepository(database.GetDB())
	userProfileRepo := repository.NewUserProfileRepository(database.GetDB())

	userService := services.NewUserService(userRepo, userProfileRepo)

	userServer := userGrpc.NewUserServer(userService) // gRPC server

	userHandler := handlers.NewUserHandler(userService)

	go startGRPCServer(userServer, cfg)

	startHTTPServer(cfg, userHandler, authServiceClient)
}

func startGRPCServer(userServer *userGrpc.UserServer, cfg *config.Config) {

	grpcPort, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Printf("Invalid port configuration: %v", err)
		return
	}
	grpcPort = grpcPort + 1000

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Printf("Failed to listen on gRPC port: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	user_service.RegisterUserServiceServer(grpcServer, userServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Failed to start gRPC server: %v", err)
		return
	}
}

func startHTTPServer(cfg *config.Config, userHandler *handlers.UserHandler, authServiceClient *clients.AuthServiceClient) {
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

	jwtMiddleware := middleware.NewJWTMiddleware(authServiceClient)

	userGroup := r.Group("/api/v1/users")
	userGroup.Use(jwtMiddleware.ValidateToken())
	{
		userGroup.GET("/", userHandler.GetUser)
		userGroup.GET("/profile", userHandler.GetUserProfile)
		userGroup.POST("/profile", userHandler.CreateUserProfile)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
