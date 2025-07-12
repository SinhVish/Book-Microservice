package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"auth-service/internal/clients"
	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/config"

	"shared/proto/user_service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.AuthReq) (*dto.AuthRes, error)
	Login(req *dto.AuthReq) (*dto.AuthRes, error)
	RefreshToken(req *dto.RefreshTokenReq) (*dto.RefreshTokenRes, error)
}

type Claims struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type authService struct {
	credentialRepo    repository.CredentialRepository
	userServiceClient *clients.UserServiceClient
	config            *config.Config
}

func NewAuthService(credentialRepo repository.CredentialRepository, userServiceClient *clients.UserServiceClient, config *config.Config) AuthService {
	return &authService{
		credentialRepo:    credentialRepo,
		userServiceClient: userServiceClient,
		config:            config,
	}
}

func (s *authService) Register(req *dto.AuthReq) (*dto.AuthRes, error) {
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

	userID, err := s.createUserInUserService(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user record: %v", err)
	}

	// Generate both tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(credential.Email, userID, credential.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &dto.AuthRes{
		ID:           credential.ID,
		Email:        credential.Email,
		Message:      "User registered successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) Login(req *dto.AuthReq) (*dto.AuthRes, error) {
	existingCredential, err := s.credentialRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %v", err)
	}
	if existingCredential == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingCredential.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get user ID from user-service
	userID, err := s.getUserIDFromUserService(existingCredential.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	// Generate both tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(existingCredential.Email, userID, existingCredential.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &dto.AuthRes{
		ID:           existingCredential.ID,
		Email:        existingCredential.Email,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) RefreshToken(req *dto.RefreshTokenReq) (*dto.RefreshTokenRes, error) {
	oldRefreshToken, err := s.credentialRepo.GetRefreshTokenByToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %v", err)
	}
	if oldRefreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}
	if oldRefreshToken.IsRevoked {
		return nil, errors.New("refresh token revoked")
	}
	if oldRefreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	credential, err := s.credentialRepo.GetByID(oldRefreshToken.CredentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %v", err)
	}

	// Get user ID from user-service
	userID, err := s.getUserIDFromUserService(credential.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	// Generate both tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(credential.Email, userID, credential.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &dto.RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// -----------------------
// -- Helper functions --
// -----------------------

func (s *authService) createUserInUserService(email string) (uint, error) {
	grpcReq := &user_service.CreateUserRequest{
		Email: email,
	}

	ctx := context.Background()
	response, err := s.userServiceClient.CreateUser(ctx, grpcReq)
	if err != nil {
		return 0, fmt.Errorf("failed to create user via gRPC: %v", err)
	}

	return uint(response.Id), nil
}

func (s *authService) getUserIDFromUserService(email string) (uint, error) {
	grpcReq := &user_service.GetUserByEmailRequest{
		Email: email,
	}

	ctx := context.Background()
	response, err := s.userServiceClient.GetUserByEmail(ctx, grpcReq)
	if err != nil {
		return 0, fmt.Errorf("failed to get user by email via gRPC: %v", err)
	}

	return uint(response.Id), nil
}

func (s *authService) generateTokenPair(email string, userID uint, credentialID uint) (string, string, int64, error) {
	// Generate access token
	accessExpirationTime := time.Now().Add(time.Duration(s.config.AccessTokenExpiryHours) * time.Hour)

	claims := &Claims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %v", err)
	}

	// Generate refresh token
	refreshTokenString := uuid.New().String()
	refreshExpirationTime := time.Now().Add(time.Duration(s.config.RefreshTokenExpiryHours) * time.Hour)

	refreshToken := &models.RefreshToken{
		CredentialID: credentialID,
		Token:        refreshTokenString,
		ExpiresAt:    refreshExpirationTime,
		IsRevoked:    false,
	}

	if err := s.credentialRepo.CreateRefreshToken(refreshToken); err != nil {
		return "", "", 0, fmt.Errorf("failed to create refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, accessExpirationTime.Unix(), nil
}
