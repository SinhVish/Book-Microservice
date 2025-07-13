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

type BookHandler struct {
	bookService services.BookService
}

func NewBookHandler(bookService services.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.CreateBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		utils.HandleError(c, utils.BadRequest("Invalid request body"))
		return
	}

	book, err := h.bookService.CreateBook(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid book ID"))
		return
	}

	book, err := h.bookService.GetBookByID(uint(id))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid book ID"))
		return
	}

	var req dto.UpdateBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		utils.HandleError(c, utils.BadRequest("Invalid request body"))
		return
	}

	book, err := h.bookService.UpdateBook(uint(id), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid book ID"))
		return
	}

	err = h.bookService.DeleteBook(uint(id))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (h *BookHandler) GetBooksByAuthor(c *gin.Context) {
	authorIDStr := c.Param("authorId")
	authorID, err := strconv.ParseUint(authorIDStr, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.BadRequest("Invalid author ID"))
		return
	}

	books, err := h.bookService.GetBooksByAuthorID(uint(authorID))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *BookHandler) SearchBooks(c *gin.Context) {
	var req dto.SearchBooksReq
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Println("Error binding query:", err)
		utils.HandleError(c, utils.BadRequest("Invalid search parameters"))
		return
	}

	result, err := h.bookService.SearchBooks(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
