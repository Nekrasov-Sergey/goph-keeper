package client

import (
	"bufio"
	"context"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

func (c *Client) Login(ctx context.Context, reader *bufio.Reader) error {
	loginPassword := readCredentials(reader)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.grpcClient.Login(ctx, &pb.LoginRequest{
		Login:    loginPassword.Login,
		Password: loginPassword.Password,
	})
	if err != nil {
		grpc.PrintError(err)
		return err
	}

	c.token = resp.Token
	return nil
}
