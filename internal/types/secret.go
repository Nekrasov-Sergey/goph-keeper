package types

import (
	"time"
)

type Secret struct {
	ID            int64      `db:"id"`
	Name          string     `db:"name"`
	Type          SecretType `db:"type"`
	EncryptedData []byte     `db:"encrypted_data"`
	Metadata      *string    `db:"metadata"`
	UserID        int64      `db:"user_id"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type SecretType string

const (
	SecretTypeUnknown       SecretType = ""
	SecretTypeLoginPassword SecretType = "LoginPassword"
	SecretTypeText          SecretType = "Text"
	SecretTypeBinary        SecretType = "Binary"
	SecretTypeBankCard      SecretType = "BankCard"
)

type SecretTypeRu string

const (
	SecretTypeUnknownRu       SecretTypeRu = ""
	SecretTypeLoginPasswordRu SecretTypeRu = "Логин и пароль"
	SecretTypeTextRu          SecretTypeRu = "Текст"
	SecretTypeBinaryRu        SecretTypeRu = "Файл"
	SecretTypeBankCardRu      SecretTypeRu = "Банковская карта"
)
