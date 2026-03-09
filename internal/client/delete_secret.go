package client

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

func (c *Client) DeleteSecret(ctx context.Context, reader *bufio.Reader) {
	fmt.Print("ID: ")
	id := readInt(reader)

	ctx, cancel := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel()

	_, err := c.grpcClient.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Id: id,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	fmt.Println("Секрет удалён")
}
