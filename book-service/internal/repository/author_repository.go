package repository

import (
	"book-service/internal/models"

	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(author *models.Author) error
	GetByID(id uint) (*models.Author, error)
	GetAll() ([]models.Author, error)
	Update(author *models.Author) error
	Delete(id uint) error
	GetByName(name string) (*models.Author, error)
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) Create(author *models.Author) error {
	return r.db.Create(author).Error
}

func (r *authorRepository) GetByID(id uint) (*models.Author, error) {
	var author models.Author
	if err := r.db.Where("id = ?", id).First(&author).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

func (r *authorRepository) GetAll() ([]models.Author, error) {
	var authors []models.Author
	if err := r.db.Find(&authors).Error; err != nil {
		return nil, err
	}
	return authors, nil
}

func (r *authorRepository) Update(author *models.Author) error {
	return r.db.Save(author).Error
}

func (r *authorRepository) Delete(id uint) error {
	return r.db.Delete(&models.Author{}, id).Error
}

func (r *authorRepository) GetByName(name string) (*models.Author, error) {
	var author models.Author
	if err := r.db.Where("name = ?", name).First(&author).Error; err != nil {
		return nil, err
	}
	return &author, nil
}
