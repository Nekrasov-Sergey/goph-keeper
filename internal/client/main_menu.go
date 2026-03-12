package client

import (
	"context"
)

func (c *Client) MainMenu(ctx context.Context) {
	for {
		cmd := selectMenu("Меню", []string{
			"Создать секрет",
			"Список секретов",
			"Получить секрет",
			"Редактировать секрет",
			"Удалить секрет",
			"Выход",
		})

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

		case "Выход":
			return
		}
	}
}
