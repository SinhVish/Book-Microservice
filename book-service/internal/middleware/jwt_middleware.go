package middleware

import (
	"log"
	"net/http"
	"strings"

	"book-service/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	config *config.Config
}

type Claims struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(config *config.Config) *JWTMiddleware {
	return &JWTMiddleware{config: config}
}

func (m *JWTMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		log.Println("authHeader", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKeyType
			}
			return []byte(m.config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid",
			})
			c.Abort()
			return
		}

		log.Println("claims", claims)

		// Set user context for downstream handlers
		c.Set("user_email", claims.Email)
		c.Set("user_id", claims.UserID)
		c.Set("user_claims", claims)

		// Continue to next handler
		c.Next()
	}
}
