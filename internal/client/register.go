// Package client реализует CLI-клиент для взаимодействия с сервером.
package client

import (
	"context"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

// Register регистрирует нового пользователя.
func (c *Client) Register(ctx context.Context) error {
	creds, err := promptCredentials()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.GRPCClient.Register(ctx, &pb.RegisterRequest{
		Login:    creds.Login,
		Password: creds.Password,
	})
	if err != nil {
		grpc.PrintError(err)
		return err
	}

	c.token = resp.Token
	return nil
}
