package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
)

func (c *Client) GetSecrets(ctx context.Context) {
	ctx, cancel := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel()

	resp, err := c.grpcClient.GetSecrets(ctx, &pb.GetSecretsRequest{})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	if len(resp.Secrets) == 0 {
		fmt.Println()
		fmt.Println("секретов нет")
		return
	}

	for _, s := range resp.Secrets {
		secretTypeRu, err := mappers.TranslateProtoSecretType(s.Type)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println()
		fmt.Println("ID:", s.Id)
		fmt.Println("Название:", s.Name)
		fmt.Println("Тип:", secretTypeRu)
	}
}
