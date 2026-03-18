// Package client реализует CLI-клиент для взаимодействия с сервером.
package client

import (
	"context"
	"fmt"
)

// MainMenu отображает главное меню работы с секретами.
func (c *Client) MainMenu(ctx context.Context) {
	for {
		fmt.Println()

		cmd, err := selectMenu("Меню", []string{
			"Создать секрет",
			"Список секретов",
			"Получить секрет",
			"Редактировать секрет",
			"Удалить секрет",
			"Назад",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch cmd {
		case "Создать секрет":
			c.CreateSecret(ctx)

		case "Список секретов":
			c.GetSecrets(ctx)

		case "Получить секрет":
			c.GetSecret(ctx)

		case "Редактировать секрет":
			c.UpdateSecret(ctx)

		case "Удалить секрет":
			c.DeleteSecret(ctx)

		case "Назад":
			return
		}
	}
}
