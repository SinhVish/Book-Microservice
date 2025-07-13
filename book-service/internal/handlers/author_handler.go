package handlers

import (
	"log"
	"net/http"
	"strconv"

	"book-service/internal/dto"
	"book-service/internal/services"

	"shared/utils"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	authorService services.AuthorService
}

func NewAuthorHandler(authorService services.AuthorService) *AuthorHandler {
	return &AuthorHandler{authorService: authorService}
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var req dto.CreateAuthorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		utils.HandleError(c, utils.BadRequest("Invalid request body"))
		return
	}

	author, err := h.authorService.CreateAuthor(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, author)
}

func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid author ID"))
		return
	}

	author, err := h.authorService.GetAuthorByID(uint(id))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, author)
}

func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	authors, err := h.authorService.GetAllAuthors()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"authors": authors})
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid author ID"))
		return
	}

	var req dto.UpdateAuthorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		utils.HandleError(c, utils.BadRequest("Invalid request body"))
		return
	}

	author, err := h.authorService.UpdateAuthor(uint(id), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, author)
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid author ID"))
		return
	}

	err = h.authorService.DeleteAuthor(uint(id))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Author deleted successfully"})
}
