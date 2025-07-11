package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"auth-service/internal/clients"
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
	Register(req *RegisterRequest) (*RegisterResponse, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	RefreshToken(req *RefreshTokenRequest) (*RefreshTokenResponse, error)
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
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewAuthService(credentialRepo repository.CredentialRepository, userServiceClient *clients.UserServiceClient, config *config.Config) AuthService {
	return &authService{
		credentialRepo:    credentialRepo,
		userServiceClient: userServiceClient,
		config:            config,
	}
}

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

	// Generate both tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(credential.Email, credential.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &RegisterResponse{
		ID:           credential.ID,
		Email:        credential.Email,
		Message:      "User registered successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) Login(req *LoginRequest) (*LoginResponse, error) {
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

	// Generate both tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(existingCredential.Email, existingCredential.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &LoginResponse{
		ID:           existingCredential.ID,
		Email:        existingCredential.Email,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) RefreshToken(req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
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

	// Generate both tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokenPair(credential.Email, credential.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// -----------------------
// -- Helper functions --
// -----------------------

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

// generateTokenPair generates both access and refresh tokens in a single operation
func (s *authService) generateTokenPair(email string, credentialID uint) (string, string, int64, error) {
	// Generate access token
	accessExpirationTime := time.Now().Add(time.Duration(s.config.AccessTokenExpiryHours) * time.Hour)

	claims := &Claims{
		Email: email,
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
