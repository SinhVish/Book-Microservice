package repository

import (
	"user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository is used to
// create a new instance of user repository with database connection
// It's important because it provides dependency injection for database operations
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create is used to
// insert a new user record into the database
// It's important because it handles user creation data persistence with proper error handling
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByEmail is used to
// retrieve a user record by email address
// It's important because the workflow is: check if user exists for authentication and user lookup operations
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID is used to
// retrieve a user record by its primary key
// It's important because it enables user lookup by ID for profile operations and user management
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update is used to
// modify an existing user record in the database
// It's important because the workflow is: update user status, email changes, or other user attributes
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete is used to
// remove a user record from the database (soft delete with GORM)
// It's important because it handles user account deletion while maintaining data integrity
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
