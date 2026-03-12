package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
)

func (c *Client) UpdateSecret(ctx context.Context) {
	id, err := promptSecretID()
	if err != nil {
		fmt.Println("ошибка ввода:", err)
		return
	}

	ctx1, cancel1 := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel1()

	resp, err := c.grpcClient.GetSecret(ctx1, &pb.GetSecretRequest{
		Id: id,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	secretType, err := mappers.ProtoSecretTypeToDomain(resp.Type)
	if err != nil {
		fmt.Println(err)
		return
	}

	updatedSecret, err := promptUpdateSecret(secretType)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx2, cancel2 := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel2()

	_, err = c.grpcClient.UpdateSecret(ctx2, &pb.UpdateSecretRequest{
		Id:       id,
		Name:     updatedSecret.Name,
		Data:     updatedSecret.Data,
		Metadata: updatedSecret.Metadata,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	fmt.Println("Секрет обновлён")
}
