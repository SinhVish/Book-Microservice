package repository

import (
	"book-service/internal/dto"
	"book-service/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *models.Book) error
	GetByID(id uint) (*models.Book, error)
	GetAll() ([]models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
	GetByAuthorID(authorID uint) ([]models.Book, error)
	SearchBooks(req dto.SearchBooksReq) ([]dto.BookWithAuthor, int64, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := r.db.Preload("Author").Where("id = ?", id).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Preload("Author").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *bookRepository) GetByAuthorID(authorID uint) ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Preload("Author").Where("author_id = ?", authorID).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) SearchBooks(req dto.SearchBooksReq) ([]dto.BookWithAuthor, int64, error) {
	var books []dto.BookWithAuthor
	var total int64

	// Build the query
	query := r.db.Model(&models.Book{}).
		Select("books.id, books.title, books.description, books.publish_year, books.isbn, books.genre, books.pages, books.price, books.author_id, authors.name as author_name, books.created_at, books.updated_at").
		Joins("JOIN authors ON books.author_id = authors.id")

	// Apply filters
	if req.AuthorName != "" {
		query = query.Where("authors.name LIKE ?", "%"+req.AuthorName+"%")
	}
	if req.BookTitle != "" {
		query = query.Where("books.title LIKE ?", "%"+req.BookTitle+"%")
	}
	if req.PublishYear > 0 {
		query = query.Where("books.publish_year = ?", req.PublishYear)
	}
	if req.Genre != "" {
		query = query.Where("books.genre LIKE ?", "%"+req.Genre+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if req.Page > 0 && req.Limit > 0 {
		offset := (req.Page - 1) * req.Limit
		query = query.Offset(offset).Limit(req.Limit)
	}

	// Execute query
	if err := query.Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}
