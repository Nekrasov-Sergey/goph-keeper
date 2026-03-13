package client

import (
	"context"
	"fmt"

	buildinfo "github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info"
)

func (c *Client) AuthMenu(ctx context.Context) {
	for {
		fmt.Println()

		cmd, err := selectMenu("Менеджер паролей GophKeeper", []string{
			"Авторизоваться",
			"Зарегистрироваться",
			"Версия и дата сборки",
			"Выход",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch cmd {
		case "Авторизоваться":
			if err := c.Login(ctx); err == nil {
				c.MainMenu(ctx)
			}

		case "Зарегистрироваться":
			if err := c.Register(ctx); err == nil {
				c.MainMenu(ctx)
			}

		case "Версия и дата сборки":
			buildinfo.Print()

		case "Выход":
			return
		}
	}
}
