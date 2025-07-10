package handlers

import (
	"log"
	"net/http"

	"auth-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest

	// ShouldBindJSON is used to bind the request body to the RegisterRequest struct
	// Meaning it will parse the request body and populate the req struct with the data
	// This is a serializer for JSON data
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		// log the error
		log.Println("Error:", err.Error())

		//! Not a good practice to check for error messages like this
		// Will improve this later
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
		})

		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login endpoint - coming soon",
	})
}

func (h *AuthHandler) Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Token validation endpoint - coming soon",
	})
}
