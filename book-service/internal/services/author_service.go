package services

import (
	"errors"

	"book-service/internal/dto"
	"book-service/internal/models"
	"book-service/internal/repository"

	"shared/utils"

	"gorm.io/gorm"
)

type AuthorService interface {
	CreateAuthor(req dto.CreateAuthorReq) (*models.Author, error)
	GetAuthorByID(id uint) (*models.Author, error)
	GetAllAuthors() ([]models.Author, error)
	UpdateAuthor(id uint, req dto.UpdateAuthorReq) (*models.Author, error)
	DeleteAuthor(id uint) error
}

type authorService struct {
	authorRepo repository.AuthorRepository
	bookRepo   repository.BookRepository
}

func NewAuthorService(authorRepo repository.AuthorRepository, bookRepo repository.BookRepository) AuthorService {
	return &authorService{
		authorRepo: authorRepo,
		bookRepo:   bookRepo,
	}
}

func (s *authorService) CreateAuthor(req dto.CreateAuthorReq) (*models.Author, error) {
	existingAuthor, err := s.authorRepo.GetByName(req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.InternalServerError("Failed to check existing author")
	}
	if existingAuthor != nil {
		return nil, utils.Conflict("Author with this name already exists")
	}

	author := &models.Author{
		Name:      req.Name,
		Bio:       req.Bio,
		BirthDate: req.BirthDate,
		Country:   req.Country,
	}

	if err := s.authorRepo.Create(author); err != nil {
		return nil, utils.InternalServerError("Failed to create author")
	}

	return author, nil
}

func (s *authorService) GetAuthorByID(id uint) (*models.Author, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("Author not found")
		}
		return nil, utils.InternalServerError("Failed to get author")
	}
	return author, nil
}

func (s *authorService) GetAllAuthors() ([]models.Author, error) {
	authors, err := s.authorRepo.GetAll()
	if err != nil {
		return nil, utils.InternalServerError("Failed to get authors")
	}
	return authors, nil
}

func (s *authorService) UpdateAuthor(id uint, req dto.UpdateAuthorReq) (*models.Author, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("Author not found")
		}
		return nil, utils.InternalServerError("Failed to get author")
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if another author with the same name exists
		existingAuthor, err := s.authorRepo.GetByName(*req.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.InternalServerError("Failed to check existing author")
		}
		if existingAuthor != nil && existingAuthor.ID != id {
			return nil, utils.Conflict("Author with this name already exists")
		}
		author.Name = *req.Name
	}
	if req.Bio != nil {
		author.Bio = *req.Bio
	}
	if req.BirthDate != nil {
		author.BirthDate = req.BirthDate
	}
	if req.Country != nil {
		author.Country = *req.Country
	}

	if err := s.authorRepo.Update(author); err != nil {
		return nil, utils.InternalServerError("Failed to update author")
	}

	return author, nil
}

func (s *authorService) DeleteAuthor(id uint) error {
	_, err := s.authorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Author not found")
		}
		return utils.InternalServerError("Failed to get author")
	}

	// Check if author has books
	books, err := s.bookRepo.GetByAuthorID(id)
	if err != nil {
		return utils.InternalServerError("Failed to check author's books")
	}
	if len(books) > 0 {
		return utils.Conflict("Cannot delete author with existing books")
	}

	if err := s.authorRepo.Delete(id); err != nil {
		return utils.InternalServerError("Failed to delete author")
	}

	return nil
}
