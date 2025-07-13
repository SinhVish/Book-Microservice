package main

import (
	"log"

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

	// Initialize repositories
	authorRepo := repository.NewAuthorRepository(database.GetDB())
	bookRepo := repository.NewBookRepository(database.GetDB())

	// Initialize services
	authorService := services.NewAuthorService(authorRepo, bookRepo)
	bookService := services.NewBookService(bookRepo, authorRepo)

	// Initialize handlers
	authorHandler := handlers.NewAuthorHandler(authorService)
	bookHandler := handlers.NewBookHandler(bookService)

	log.Printf("Starting HTTP server...")
	startHTTPServer(cfg, authorHandler, bookHandler)
}

func startHTTPServer(cfg *config.Config, authorHandler *handlers.AuthorHandler, bookHandler *handlers.BookHandler) {
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

	// Initialize JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg)

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Author routes (protected)
		authors := v1.Group("/authors")
		authors.Use(jwtMiddleware.ValidateToken())
		{
			authors.POST("/", authorHandler.CreateAuthor)
			authors.GET("/", authorHandler.GetAuthors)
			authors.GET("/:id", authorHandler.GetAuthor)
			authors.PUT("/:id", authorHandler.UpdateAuthor)
			authors.DELETE("/:id", authorHandler.DeleteAuthor)
		}

		// Book routes (protected)
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

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
