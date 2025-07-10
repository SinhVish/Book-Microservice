package repository

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type CredentialRepository interface {
	Create(credential *models.Credential) error
	GetByEmail(email string) (*models.Credential, error)
	GetByID(id uint) (*models.Credential, error)
	Update(credential *models.Credential) error
	Delete(id uint) error
}

type credentialRepository struct {
	db *gorm.DB
}

// NewCredentialRepository is used to
// create a new instance of credential repository with database connection
// It's important because it provides dependency injection for database operations
func NewCredentialRepository(db *gorm.DB) CredentialRepository {
	return &credentialRepository{db: db}
}

// Create is used to
// insert a new credential record into the database
// It's important because it handles user registration data persistence with proper error handling
func (r *credentialRepository) Create(credential *models.Credential) error {
	return r.db.Create(credential).Error
}

// GetByEmail is used to
// retrieve a credential record by email address
// It's important because the workflow is: check if email exists during registration and find user credentials during login
func (r *credentialRepository) GetByEmail(email string) (*models.Credential, error) {
	var credential models.Credential
	err := r.db.Where("email = ?", email).First(&credential).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// GetByID is used to
// retrieve a credential record by its primary key
// It's important because it enables user lookup by ID for authentication operations and profile updates
func (r *credentialRepository) GetByID(id uint) (*models.Credential, error) {
	var credential models.Credential
	err := r.db.First(&credential, id).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// Update is used to
// modify an existing credential record in the database
// It's important because the workflow is: update last login time, change password, or modify account status
func (r *credentialRepository) Update(credential *models.Credential) error {
	return r.db.Save(credential).Error
}

// Delete is used to
// remove a credential record from the database (soft delete with GORM)
// It's important because it handles user account deletion while maintaining data integrity
func (r *credentialRepository) Delete(id uint) error {
	return r.db.Delete(&models.Credential{}, id).Error
}
