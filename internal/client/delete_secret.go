package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

func (c *Client) DeleteSecret(ctx context.Context) {
	id, err := promptSecretID()
	if err != nil {
		fmt.Println("ошибка ввода:", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel()

	_, err = c.GRPCClient.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Id: id,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	fmt.Println()
	fmt.Println("Секрет удалён")
}
