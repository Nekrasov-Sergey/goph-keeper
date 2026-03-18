// Package errcodes содержит коды ошибок приложения.
package errcodes

import (
	"github.com/pkg/errors"
)

var (
	// ErrLoginAlreadyExists — ошибка: логин уже занят.
	ErrLoginAlreadyExists = errors.New("логин уже занят")
	// ErrInvalidCredentials — ошибка: неверная пара логин/пароль.
	ErrInvalidCredentials = errors.New("неверная пара логин/пароль")
	// ErrUserNotFound — ошибка: пользователь не найден.
	ErrUserNotFound = errors.New("пользователь не найден")
	// ErrSecretNameAlreadyExists — ошибка: имя секрета уже занято.
	ErrSecretNameAlreadyExists = errors.New("имя секрета уже занято")
	// ErrSecretNotFound — ошибка: секрет не найден.
	ErrSecretNotFound = errors.New("секрет не найден")
)
