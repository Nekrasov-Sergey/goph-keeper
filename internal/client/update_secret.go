package client

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

func (c *Client) UpdateSecret(ctx context.Context, reader *bufio.Reader) {
	fmt.Print("ID: ")
	id := readInt(reader)

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

	updatedSecret := readUpdateSecretInput(ctx, reader, secretType)
	if updatedSecret == nil {
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
}

func readUpdateSecretInput(ctx context.Context, reader *bufio.Reader, secretType types.SecretType) *types.UpdatedSecret {
	updatedSecret := &types.UpdatedSecret{}

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Println()
		fmt.Println("Какой параметр обновить?")
		fmt.Println("1. Название")
		fmt.Println("2. Основные данные")
		fmt.Println("3. Дополнительные данные")
		fmt.Println("0. Выйти из программы")
		fmt.Print("> ")

		cmd := readLine(reader)

		switch cmd {
		case "1":
			fmt.Print("Название: ")
			updatedSecret.Name = utils.Ptr(readLine(reader))
			break loop

		case "2":
			switch secretType {
			case types.SecretTypeLoginPassword:
				data, ok := readCreateLoginPassword(reader)
				if !ok {
					continue
				}

				updatedSecret.Data = data
				break loop

			case types.SecretTypeText:
				data, ok := readCreateText(reader)
				if !ok {
					continue
				}

				updatedSecret.Data = data
				break loop

			case types.SecretTypeBinary:
				data, ok := readCreateBinary(reader)
				if !ok {
					continue
				}

				updatedSecret.Data = data
				break loop

			case types.SecretTypeBankCard:
				data, ok := readBankCard(reader)
				if !ok {
					continue
				}

				updatedSecret.Data = data
				break loop
			}

			break loop

		case "3":
			fmt.Print("Дополнительные данные: ")
			updatedSecret.Metadata = utils.Ptr(readLine(reader))
			break loop

		case "0":
			fmt.Println("Выход")
			return nil

		default:
			fmt.Println("Неизвестная команда")
		}
	}

	return updatedSecret
}
