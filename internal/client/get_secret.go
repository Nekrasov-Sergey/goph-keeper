package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

func (c *Client) GetSecret(ctx context.Context, reader *bufio.Reader) {
	fmt.Print("ID: ")
	id := readInt(reader)

	ctx, cancel := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel()

	resp, err := c.grpcClient.GetSecret(ctx, &pb.GetSecretRequest{
		Id: id,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}

	fmt.Println("Название:", resp.Name)

	secretTypeRu, err := mappers.TranslateProtoSecretType(resp.Type)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Тип:", secretTypeRu)

	secretType, err := mappers.ProtoSecretTypeToDomain(resp.Type)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch secretType {
	case types.SecretTypeLoginPassword:
		var loginPassword types.LoginPassword
		if err := json.Unmarshal(resp.Data, &loginPassword); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Логин:", loginPassword.Login)
		fmt.Println("Пароль:", loginPassword.Password)

	case types.SecretTypeText:
		var text string
		if err := json.Unmarshal(resp.Data, &text); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(text)

	case types.SecretTypeBinary:
		fmt.Print("Куда сохранить файл: ")
		path := readLine(reader)

		if err := os.WriteFile(path, resp.Data, 0600); err != nil {
			fmt.Println("Ошибка сохранения файла:", err)
			return
		}

		fmt.Println("Файл успешно сохранён")

	case types.SecretTypeBankCard:
		var bankCard types.BankCard
		if err := json.Unmarshal(resp.Data, &bankCard); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Номер карты:", bankCard.Number)
		fmt.Println("Владелец карты:", bankCard.Holder)
		fmt.Println("Срок действия (MM/YY):", bankCard.Expiry)
		fmt.Println("CVV код:", bankCard.CVV)
	}

	if resp.Metadata != nil && utils.Deref(resp.Metadata) != "" {
		fmt.Println("Дополнительная информация:", utils.Deref(resp.Metadata))
	}
}
