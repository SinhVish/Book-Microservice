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

func (c *UserServiceClient) CreateUser(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return response, nil
}

func (c *UserServiceClient) GetUserByEmail(ctx context.Context, req *user_service.GetUserByEmailRequest) (*user_service.GetUserByEmailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := c.client.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return response, nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
