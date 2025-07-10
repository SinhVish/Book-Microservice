package services

import (
	"context"
	"errors"
	"fmt"

	"auth-service/internal/clients"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/config"

	"shared/proto/user_service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *RegisterRequest) (*RegisterResponse, error)
}

type authService struct {
	credentialRepo    repository.CredentialRepository
	userServiceClient *clients.UserServiceClient
	config            *config.Config
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// NewAuthService is used to
// create a new instance of authentication service with injected dependencies
// It's important because it provides dependency injection for repository, gRPC client, and config, enabling clean architecture and testability
func NewAuthService(credentialRepo repository.CredentialRepository, userServiceClient *clients.UserServiceClient, config *config.Config) AuthService {
	return &authService{
		credentialRepo:    credentialRepo,
		userServiceClient: userServiceClient,
		config:            config,
	}
}

// Register is used to
// handle the complete user registration business logic including validation, password hashing, and gRPC microservice coordination
// It's important because the workflow is: check email uniqueness, hash password, create credentials, call user service via gRPC to create basic user record, and return response
func (s *authService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	existingCredential, err := s.credentialRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing email: %v", err)
	}
	if existingCredential != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	credential := &models.Credential{
		Email:    req.Email,
		Password: string(hashedPassword),
		IsActive: true,
	}

	if err := s.credentialRepo.Create(credential); err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	if err := s.createUserInUserService(req.Email); err != nil {
		return nil, fmt.Errorf("failed to create user record: %v", err)
	}

	return &RegisterResponse{
		ID:      credential.ID,
		Email:   credential.Email,
		Message: "User registered successfully",
	}, nil
}

// createUserInUserService is used to
// communicate with the user service via gRPC to create basic user record during registration
// It's important because the workflow is: prepare gRPC request with only email, call user service using gRPC client to create user record, and handle response for microservice coordination
func (s *authService) createUserInUserService(email string) error {
	grpcReq := &user_service.CreateUserRequest{
		Email: email,
	}

	ctx := context.Background()
	_, err := s.userServiceClient.CreateUser(ctx, grpcReq)
	if err != nil {
		return fmt.Errorf("failed to create user via gRPC: %v", err)
	}

	return nil
}
