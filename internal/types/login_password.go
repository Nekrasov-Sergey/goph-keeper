// Package types содержит доменные типы данных приложения.
package types

// LoginPassword представляет пару логин/пароль для хранения в секрете.
type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
