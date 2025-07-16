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

	// Initialize cached auth service client
	authServiceClient, err := clients.NewCachedAuthClient(cfg.AuthServiceURL, cfg.RedisURL, cfg.GetL1CacheTTL(), cfg.GetL2CacheTTL())
	if err != nil {
		log.Fatal("Failed to create auth service client:", err)
	}
	defer authServiceClient.Close()

	userRepo := repository.NewUserRepository(database.GetDB())
	userProfileRepo := repository.NewUserProfileRepository(database.GetDB())

	userService := services.NewUserService(userRepo, userProfileRepo)

	userServer := userGrpc.NewUserServer(userService) // gRPC server

	userHandler := handlers.NewUserHandler(userService)
	cacheHandler := handlers.NewCacheHandler(authServiceClient)

	log.Printf("Starting gRPC server in goroutine...")
	go startGRPCServer(userServer, cfg)

	log.Printf("Starting HTTP server...")
	startHTTPServer(cfg, userHandler, cacheHandler, authServiceClient)
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

func startHTTPServer(cfg *config.Config, userHandler *handlers.UserHandler, cacheHandler *handlers.CacheHandler, authServiceClient *clients.CachedAuthClient) {
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

	// Cache monitoring endpoints
	cacheGroup := r.Group("/cache")
	{
		cacheGroup.GET("/stats", cacheHandler.GetCacheStats)
		cacheGroup.GET("/metrics", cacheHandler.GetCacheMetrics)
		cacheGroup.POST("/clear", cacheHandler.ClearCache)
	}

	jwtMiddleware := middleware.NewJWTMiddleware(authServiceClient)

	userGroup := r.Group("/api/v1/users")
	userGroup.Use(jwtMiddleware.ValidateToken())
	{
		userGroup.GET("/", userHandler.GetUser)
		userGroup.GET("/profile", userHandler.GetUserProfile)
		userGroup.POST("/profile", userHandler.CreateUserProfile)
	}

	log.Printf("HTTP server starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("Auth Service URL: %s", cfg.AuthServiceURL)
	log.Printf("Redis URL: %s", cfg.RedisURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
