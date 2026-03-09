package client

import (
	"bufio"
	"context"
	"fmt"
)

func (c *Client) Menu(ctx context.Context, reader *bufio.Reader) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Println()
		fmt.Println("1. Создать секрет")
		fmt.Println("2. Получить список секретов")
		fmt.Println("3. Получить секрет")
		fmt.Println("4. Редактировать секрет")
		fmt.Println("5. Удалить секрет")
		fmt.Println("0. Выйти из программы")
		fmt.Print("> ")

		cmd := readLine(reader)

		switch cmd {

		case "1":
			c.CreateSecret(ctx, reader)

		case "2":
			c.GetSecrets(ctx)

		case "3":
			c.GetSecret(ctx, reader)

		case "4":
			c.UpdateSecret(ctx, reader)

		case "5":
			c.DeleteSecret(ctx, reader)

		case "0":
			return

		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
