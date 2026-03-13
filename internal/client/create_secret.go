package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
)

func (c *Client) CreateSecret(ctx context.Context) {
	input, err := promptCreateSecret()
	if err != nil {
		fmt.Println(err)
		return
	}

	secretType, err := mappers.DomainSecretTypeToProto(input.Type)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel()

	_, err = c.grpcClient.CreateSecret(ctx, &pb.CreateSecretRequest{
		Name:     input.Name,
		Type:     secretType,
		Data:     input.Data,
		Metadata: input.Metadata,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	fmt.Println()
	fmt.Println("Секрет создан")
}
