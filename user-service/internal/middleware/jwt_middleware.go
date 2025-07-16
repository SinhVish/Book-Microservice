package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"user-service/internal/clients"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	authClient *clients.CachedAuthClient
}

func NewJWTMiddleware(authClient *clients.CachedAuthClient) *JWTMiddleware {
	return &JWTMiddleware{authClient: authClient}
}

func (m *JWTMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Println("authHeader", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Call auth-service to validate token via gRPC
		ctx := context.Background()
		response, err := m.authClient.ValidateToken(ctx, tokenString)
		if err != nil {
			log.Printf("Failed to call auth service: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
			c.Abort()
			return
		}

		if !response.IsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": response.ErrorMessage})
			c.Abort()
			return
		}

		log.Printf("Token validated successfully for user: %s", response.Claims.Email)

		// Set user information in context
		c.Set("user_email", response.Claims.Email)
		c.Set("user_id", uint(response.Claims.UserId))
		c.Set("user_claims", response.Claims)

		c.Next()
	}
}
