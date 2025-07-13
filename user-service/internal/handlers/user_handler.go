package handlers

import (
	"log"
	"net/http"

	"user-service/internal/dto"
	"user-service/internal/services"

	"shared/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateUserProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		utils.HandleError(c, utils.BadRequest("Invalid request body"))
		return
	}

	user_profile, err := h.userService.CreateUserProfile(userID, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user_profile)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	profile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		// Just like Django - one line handles all custom errors!
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}
