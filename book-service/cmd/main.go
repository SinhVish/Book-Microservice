package main

import (
	"log"

	"book-service/internal/clients"
	"book-service/internal/handlers"
	"book-service/internal/middleware"
	"book-service/internal/repository"
	"book-service/internal/services"
	"book-service/pkg/config"
	"book-service/pkg/database"

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

	// Initialize cached auth service client
	authServiceClient, err := clients.NewCachedAuthClient(cfg.AuthServiceURL, cfg.RedisURL, cfg.GetL1CacheTTL(), cfg.GetL2CacheTTL())
	if err != nil {
		log.Fatal("Failed to create auth service client:", err)
	}
	defer authServiceClient.Close()

	// Initialize repositories
	authorRepo := repository.NewAuthorRepository(database.GetDB())
	bookRepo := repository.NewBookRepository(database.GetDB())

	// Initialize services
	authorService := services.NewAuthorService(authorRepo, bookRepo)
	bookService := services.NewBookService(bookRepo, authorRepo)

	// Initialize handlers
	authorHandler := handlers.NewAuthorHandler(authorService)
	bookHandler := handlers.NewBookHandler(bookService)
	cacheHandler := handlers.NewCacheHandler(authServiceClient)

	log.Printf("Starting HTTP server...")
	startHTTPServer(cfg, authorHandler, bookHandler, cacheHandler, authServiceClient)
}

func startHTTPServer(cfg *config.Config, authorHandler *handlers.AuthorHandler, bookHandler *handlers.BookHandler, cacheHandler *handlers.CacheHandler, authServiceClient *clients.CachedAuthClient) {
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "book-service",
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

	v1 := r.Group("/api/v1")
	{
		authors := v1.Group("/authors")
		authors.Use(jwtMiddleware.ValidateToken())
		{
			authors.POST("/", authorHandler.CreateAuthor)
			authors.GET("/", authorHandler.GetAuthors)
			authors.GET("/:id", authorHandler.GetAuthor)
			authors.PUT("/:id", authorHandler.UpdateAuthor)
			authors.DELETE("/:id", authorHandler.DeleteAuthor)
		}

		books := v1.Group("/books")
		books.Use(jwtMiddleware.ValidateToken())
		{
			books.POST("/", bookHandler.CreateBook)
			books.GET("/", bookHandler.GetBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
			books.GET("/author/:authorId", bookHandler.GetBooksByAuthor)
			books.GET("/search", bookHandler.SearchBooks)
		}
	}

	log.Printf("HTTP server starting on port %s", cfg.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	log.Printf("Auth Service URL: %s", cfg.AuthServiceURL)
	log.Printf("Redis URL: %s", cfg.RedisURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
