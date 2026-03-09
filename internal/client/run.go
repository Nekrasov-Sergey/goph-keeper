package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	buildinfo "github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info"
)

func (c *Client) Run(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("Менеджер паролей GophKeeper")

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Println()
		fmt.Println("1. Авторизоваться")
		fmt.Println("2. Зарегистрироваться")
		fmt.Println("9. Версия и дата сборки")
		fmt.Println("0. Выйти из программы")
		fmt.Print("> ")

		cmd := readLine(reader)

		switch cmd {
		case "1":
			if err := c.Login(ctx, reader); err != nil {
				continue
			}
			c.Menu(ctx, reader)
			return

		case "2":
			if err := c.Register(ctx, reader); err != nil {
				continue
			}
			c.Menu(ctx, reader)
			return

		case "9":
			buildinfo.Print()

		case "0":
			return

		default:
			fmt.Println("Неизвестная команда")
		}
	}
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func readInt(reader *bufio.Reader) int64 {
	v := readLine(reader)
	id, _ := strconv.ParseInt(v, 10, 64)
	return id
}

func readCredentials(reader *bufio.Reader) types.LoginPassword {
	fmt.Print("Логин: ")
	login := readLine(reader)

	fmt.Print("Пароль: ")
	password := readLine(reader)

	return types.LoginPassword{
		Login:    login,
		Password: password,
	}
}
