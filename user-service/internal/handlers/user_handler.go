package handlers

import (
	"log"
	"net/http"

	"user-service/internal/dto"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	log.Println("GetUser")
	userID := c.GetUint("user_id")
	log.Println("userID", userID)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUserProfile(c *gin.Context) {
	// Extract user_id from context (set by JWT middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		log.Println("user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		log.Println("Invalid user_id type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	var req dto.CreateUserProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_profile, err := h.userService.CreateUserProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user_profile)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Extract user_id from context (set by JWT middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		log.Println("user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		log.Println("Invalid user_id type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	profile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
