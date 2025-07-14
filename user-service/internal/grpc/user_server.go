package grpc

import (
	"context"
	"log"

	"user-service/internal/services"

	"shared/proto/user_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	user_service.UnimplementedUserServiceServer
	userService services.UserService
}

func NewUserServer(userService services.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	log.Printf("Received CreateUser request for email: %s", req.Email)

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := s.userService.CreateUser(req.Email)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		if err.Error() == "user with this email already exists" {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	log.Printf("Successfully created user with ID: %d", user.ID)

	return &user_service.CreateUserResponse{
		Id:      uint32(user.ID),
		Email:   user.Email,
		Message: "User created successfully",
	}, nil
}

func (s *UserServer) GetUserByEmail(ctx context.Context, req *user_service.GetUserByEmailRequest) (*user_service.GetUserByEmailResponse, error) {
	log.Printf("Received GetUserByEmail request for email: %s", req.Email)

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := s.userService.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Failed to get user by email: %v", err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	log.Printf("Successfully retrieved user with ID: %d", user.ID)

	return &user_service.GetUserByEmailResponse{
		Id:     uint32(user.ID),
		Email:  user.Email,
		Status: user.Status,
	}, nil
}
