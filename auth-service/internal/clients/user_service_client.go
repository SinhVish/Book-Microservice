package clients

import (
	"context"
	"fmt"
	"time"

	"shared/proto/user_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	conn   *grpc.ClientConn
	client user_service.UserServiceClient
}

// NewUserServiceClient is used to
// create a new gRPC client connection to the user service
// It's important because it establishes the gRPC connection for inter-service communication
func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	client := user_service.NewUserServiceClient(conn)

	return &UserServiceClient{
		conn:   conn,
		client: client,
	}, nil
}

// CreateUser is used to
// call the user service gRPC method to create a new user profile
// It's important because the workflow is: prepare gRPC request, call user service, and handle response for microservice coordination
func (c *UserServiceClient) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return response, nil
}

// GetUserByEmail is used to
// call the user service gRPC method to retrieve user information by email
// It's important because it enables user lookup for authentication and validation purposes
func (c *UserServiceClient) GetUserByEmail(ctx context.Context, req *user_service.GetUserByEmailRequest) (*user_service.GetUserByEmailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := c.client.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return response, nil
}

// Close is used to
// close the gRPC connection to the user service
// It's important because it properly releases network resources and prevents connection leaks
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
