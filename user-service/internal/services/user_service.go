package services

import (
	"errors"

	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"

	"shared/utils"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUserProfile(userID uint, req dto.CreateUserProfileReq) (*models.UserProfile, error)
	GetUserProfile(userID uint) (*models.UserProfile, error)
	CreateUser(email string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
}

type userService struct {
	userRepo        repository.UserRepository
	userProfileRepo repository.UserProfileRepository
}

func NewUserService(userRepo repository.UserRepository, userProfileRepo repository.UserProfileRepository) UserService {
	return &userService{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
	}
}

func (s *userService) GetUserProfile(userID uint) (*models.UserProfile, error) {
	profile, err := s.userProfileRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("User profile not found")
		}
		return nil, utils.InternalServerError("Failed to get user profile")
	}
	return profile, nil
}

func (s *userService) CreateUserProfile(userID uint, req dto.CreateUserProfileReq) (*models.UserProfile, error) {
	existingProfile, err := s.userProfileRepo.GetByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.InternalServerError("Failed to check existing profile")
	}
	if existingProfile != nil {
		return nil, utils.Conflict("User profile already exists")
	}

	profile := &models.UserProfile{
		ID:          userID,
		FirstName:   &req.FirstName,
		LastName:    &req.LastName,
		Phone:       &req.Phone,
		DateOfBirth: &req.DateOfBirth,
		Gender:      &req.Gender,
		Bio:         &req.Bio,
		Address:     &req.Address,
	}

	if err := s.userProfileRepo.Create(profile); err != nil {
		return nil, utils.InternalServerError("Failed to create user profile")
	}

	return profile, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("User not found")
		}
		return nil, utils.InternalServerError("Failed to get user by email")
	}
	return user, nil
}

func (s *userService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("User not found")
		}
		return nil, utils.InternalServerError("Failed to get user by id")
	}
	return user, nil
}

func (s *userService) CreateUser(email string) (*models.User, error) {
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.InternalServerError("Failed to check existing email")
	}
	if existingUser != nil {
		return nil, utils.Conflict("User with this email already exists")
	}

	user := &models.User{
		Email:  email,
		Status: "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, utils.InternalServerError("Failed to create user")
	}

	return user, nil
}
