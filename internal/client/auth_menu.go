package client

import (
	"context"

	"github.com/manifoldco/promptui"

	buildinfo "github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info"
)

func (c *Client) AuthMenu(ctx context.Context) {
	for {
		cmd := selectMenu("Менеджер паролей GophKeeper", []string{
			"Авторизоваться",
			"Зарегистрироваться",
			"Версия и дата сборки",
			"Выход",
		})

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

func selectMenu(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return ""
	}

	return result
}
