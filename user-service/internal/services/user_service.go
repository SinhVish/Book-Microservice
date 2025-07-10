package services

import (
	"errors"
	"fmt"

	"user-service/internal/models"
	"user-service/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(email string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService is used to
// create a new instance of user service with injected dependencies
// It's important because it provides dependency injection for repository, enabling clean architecture and testability
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser is used to
// handle user creation business logic including validation and default values
// It's important because the workflow is: check email uniqueness, create user with default status, and return user record
func (s *userService) CreateUser(email string) (*models.User, error) {
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing email: %v", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user := &models.User{
		Email:  email,
		Status: "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// GetUserByEmail is used to
// retrieve user information by email address
// It's important because it enables user lookup for authentication and profile operations
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}
	return user, nil
}

// GetUserByID is used to
// retrieve user information by ID
// It's important because it enables user lookup for profile operations and user management
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}
	return user, nil
}
