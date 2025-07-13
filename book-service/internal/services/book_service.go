package services

import (
	"errors"
	"math"

	"book-service/internal/dto"
	"book-service/internal/models"
	"book-service/internal/repository"

	"shared/utils"

	"gorm.io/gorm"
)

type BookService interface {
	CreateBook(req dto.CreateBookReq) (*models.Book, error)
	GetBookByID(id uint) (*models.Book, error)
	GetAllBooks() ([]models.Book, error)
	UpdateBook(id uint, req dto.UpdateBookReq) (*models.Book, error)
	DeleteBook(id uint) error
	GetBooksByAuthorID(authorID uint) ([]models.Book, error)
	SearchBooks(req dto.SearchBooksReq) (*dto.SearchBooksRes, error)
}

type bookService struct {
	bookRepo   repository.BookRepository
	authorRepo repository.AuthorRepository
}

func NewBookService(bookRepo repository.BookRepository, authorRepo repository.AuthorRepository) BookService {
	return &bookService{
		bookRepo:   bookRepo,
		authorRepo: authorRepo,
	}
}

func (s *bookService) CreateBook(req dto.CreateBookReq) (*models.Book, error) {
	// Check if author exists
	_, err := s.authorRepo.GetByID(req.AuthorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("Author not found")
		}
		return nil, utils.InternalServerError("Failed to check author")
	}

	book := &models.Book{
		Title:       req.Title,
		Description: req.Description,
		PublishYear: req.PublishYear,
		ISBN:        req.ISBN,
		Genre:       req.Genre,
		Pages:       req.Pages,
		Price:       req.Price,
		AuthorID:    req.AuthorID,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, utils.InternalServerError("Failed to create book")
	}

	return book, nil
}

func (s *bookService) GetBookByID(id uint) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("Book not found")
		}
		return nil, utils.InternalServerError("Failed to get book")
	}
	return book, nil
}

func (s *bookService) GetAllBooks() ([]models.Book, error) {
	books, err := s.bookRepo.GetAll()
	if err != nil {
		return nil, utils.InternalServerError("Failed to get books")
	}
	return books, nil
}

func (s *bookService) UpdateBook(id uint, req dto.UpdateBookReq) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("Book not found")
		}
		return nil, utils.InternalServerError("Failed to get book")
	}

	// Update fields if provided
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.PublishYear != nil {
		book.PublishYear = *req.PublishYear
	}
	if req.ISBN != nil {
		book.ISBN = *req.ISBN
	}
	if req.Genre != nil {
		book.Genre = *req.Genre
	}
	if req.Pages != nil {
		book.Pages = *req.Pages
	}
	if req.Price != nil {
		book.Price = *req.Price
	}
	if req.AuthorID != nil {
		// Check if author exists
		_, err := s.authorRepo.GetByID(*req.AuthorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utils.NotFound("Author not found")
			}
			return nil, utils.InternalServerError("Failed to check author")
		}
		book.AuthorID = *req.AuthorID
	}

	if err := s.bookRepo.Update(book); err != nil {
		return nil, utils.InternalServerError("Failed to update book")
	}

	return book, nil
}

func (s *bookService) DeleteBook(id uint) error {
	_, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Book not found")
		}
		return utils.InternalServerError("Failed to get book")
	}

	if err := s.bookRepo.Delete(id); err != nil {
		return utils.InternalServerError("Failed to delete book")
	}

	return nil
}

func (s *bookService) GetBooksByAuthorID(authorID uint) ([]models.Book, error) {
	// Check if author exists
	_, err := s.authorRepo.GetByID(authorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("Author not found")
		}
		return nil, utils.InternalServerError("Failed to check author")
	}

	books, err := s.bookRepo.GetByAuthorID(authorID)
	if err != nil {
		return nil, utils.InternalServerError("Failed to get books by author")
	}

	return books, nil
}

func (s *bookService) SearchBooks(req dto.SearchBooksReq) (*dto.SearchBooksRes, error) {
	// Set default pagination values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	books, total, err := s.bookRepo.SearchBooks(req)
	if err != nil {
		return nil, utils.InternalServerError("Failed to search books")
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &dto.SearchBooksRes{
		Books:      books,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}
