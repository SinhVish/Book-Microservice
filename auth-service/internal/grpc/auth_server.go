package grpc

import (
	"context"
	"log"
	"time"

	"auth-service/pkg/config"

	"shared/proto/auth_service"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthServer struct {
	auth_service.UnimplementedAuthServiceServer
	config *config.Config
}

func NewAuthServer(config *config.Config) *AuthServer {
	return &AuthServer{
		config: config,
	}
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *auth_service.ValidateTokenRequest) (*auth_service.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &auth_service.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: "Token is required",
		}, nil
	}

	// Parse and validate the JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return &auth_service.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: "Invalid token: " + err.Error(),
		}, nil
	}

	if !token.Valid {
		return &auth_service.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: "Token is not valid",
		}, nil
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return &auth_service.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: "Token has expired",
		}, nil
	}

	log.Printf("Token validated successfully for user: %s", claims.Email)

	// Return successful validation with claims
	return &auth_service.ValidateTokenResponse{
		IsValid:      true,
		ErrorMessage: "",
		Claims: &auth_service.UserClaims{
			Email:     claims.Email,
			UserId:    uint32(claims.UserID),
			Issuer:    claims.Issuer,
			Subject:   claims.Subject,
			ExpiresAt: claims.ExpiresAt.Unix(),
			IssuedAt:  claims.IssuedAt.Unix(),
		},
	}, nil
}
