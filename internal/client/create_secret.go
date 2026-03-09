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

func (c *Client) CreateSecret(ctx context.Context, reader *bufio.Reader) {
	secretInput := readCreateSecretInput(ctx, reader)
	if secretInput == nil {
		return
	}

	secretType, err := mappers.DomainSecretTypeToProto(secretInput.Type)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.AuthContext(ctx), 30*time.Second)
	defer cancel()

	_, err = c.grpcClient.CreateSecret(ctx, &pb.CreateSecretRequest{
		Name:     secretInput.Name,
		Type:     secretType,
		Data:     secretInput.Data,
		Metadata: secretInput.Metadata,
	})
	if err != nil {
		grpc.PrintError(err)
		return
	}
}

func readCreateSecretInput(ctx context.Context, reader *bufio.Reader) *types.SecretInput {
	secretInput := &types.SecretInput{}

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Println()
		fmt.Println("Тип секрета")
		fmt.Println("1.", types.SecretTypeLoginPasswordRu)
		fmt.Println("2.", types.SecretTypeTextRu)
		fmt.Println("3.", types.SecretTypeBinaryRu)
		fmt.Println("4.", types.SecretTypeBankCardRu)
		fmt.Println("0. Выйти из программы")
		fmt.Print("> ")

		cmd := readLine(reader)

		switch cmd {

		case "1":
			secretInput.Type = types.SecretTypeLoginPassword

			data, ok := readCreateLoginPassword(reader)
			if !ok {
				continue
			}

			secretInput.Data = data
			break loop

		case "2":
			secretInput.Type = types.SecretTypeText

			data, ok := readCreateText(reader)
			if !ok {
				continue
			}

			secretInput.Data = data
			break loop

		case "3":
			secretInput.Type = types.SecretTypeBinary

			data, ok := readCreateBinary(reader)
			if !ok {
				continue
			}

			secretInput.Data = data
			break loop

		case "4":
			secretInput.Type = types.SecretTypeBankCard

			data, ok := readBankCard(reader)
			if !ok {
				continue
			}

			secretInput.Data = data
			break loop

		case "0":
			fmt.Println("Выход")
			return nil

		default:
			fmt.Println("Неизвестная команда")
		}
	}

	fmt.Println()
	fmt.Print("Название секрета: ")
	secretInput.Name = readLine(reader)

	secretInput.Metadata = readMetadataInput(ctx, reader)

	return secretInput
}

func readCreateLoginPassword(reader *bufio.Reader) (data []byte, ok bool) {
	loginPassword := readCredentials(reader)

	data, err := json.Marshal(loginPassword)
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return nil, false
	}

	return data, true
}

func readCreateText(reader *bufio.Reader) (data []byte, ok bool) {
	fmt.Print("Текст: ")
	text := readLine(reader)

	data, err := json.Marshal(text)
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return nil, false
	}

	return data, true
}

func readCreateBinary(reader *bufio.Reader) (data []byte, ok bool) {
	fmt.Print("Путь к файлу: ")
	path := readLine(reader)

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("Ошибка получения информации о файле:", err)
		return nil, false
	}

	if info.Size() > 10*1024*1024 {
		fmt.Println("Файл слишком большой (максимум 10MB)")
		return nil, false
	}

	data, err = os.ReadFile(path)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return nil, false
	}

	return data, true
}

func readBankCard(reader *bufio.Reader) (data []byte, ok bool) {
	fmt.Print("Номер карты: ")
	number := readLine(reader)

	fmt.Print("Владелец карты: ")
	holder := readLine(reader)

	fmt.Print("Срок действия (MM/YY): ")
	expiry := readLine(reader)

	fmt.Print("CVV код: ")
	cvv := readLine(reader)

	bankCard := types.BankCard{
		Number: number,
		Holder: holder,
		Expiry: expiry,
		CVV:    cvv,
	}

	data, err := json.Marshal(bankCard)
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return nil, false
	}

	return data, true
}

func readMetadataInput(ctx context.Context, reader *bufio.Reader) *string {
	var metadata *string

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Println()
		fmt.Println("Ввести дополнительную информацию?")
		fmt.Println("1. Да")
		fmt.Println("2. Нет")
		fmt.Println("0. Выйти из программы")
		fmt.Print("> ")

		cmd := readLine(reader)

		switch cmd {
		case "1":
			fmt.Print("Дополнительная информация: ")
			metadata = utils.Ptr(readLine(reader))
			break loop

		case "2":
			metadata = nil
			break loop

		case "0":
			fmt.Println("Выход")
			return nil

		default:
			fmt.Println("Неизвестная команда")
		}
	}

	return metadata
}
