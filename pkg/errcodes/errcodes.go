package errcodes

import (
	"github.com/pkg/errors"
)

var (
	ErrLoginAlreadyExists      = errors.New("логин уже занят")
	ErrInvalidCredentials      = errors.New("неверная пара логин/пароль")
	ErrUserNotFound            = errors.New("пользователь не найден")
	ErrSecretNameAlreadyExists = errors.New("имя секрета уже занято")
	ErrSecretNotFound          = errors.New("секрет не найден")
)
