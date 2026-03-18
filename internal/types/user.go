// Package types содержит доменные типы данных приложения.
package types

// User представляет пользователя системы.
type User struct {
	ID               int64  `db:"id"`
	Login            string `db:"login"`
	Password         string `db:"password_hash"`
	EncryptedUserKey []byte `db:"encrypted_user_key"`
}
