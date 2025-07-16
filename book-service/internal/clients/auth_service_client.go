package clients

import (
	"context"
	"fmt"
	"time"

	"shared/proto/auth_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	conn   *grpc.ClientConn
	client auth_service.AuthServiceClient
}

func NewAuthServiceClient(address string) (*AuthServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}

	client := auth_service.NewAuthServiceClient(conn)

	return &AuthServiceClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *AuthServiceClient) ValidateToken(ctx context.Context, token string) (*auth_service.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &auth_service.ValidateTokenRequest{
		Token: token,
	}

	response, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	return response, nil
}

func (c *AuthServiceClient) Close() error {
	return c.conn.Close()
}
