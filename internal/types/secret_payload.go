// Package types содержит доменные типы данных приложения.
package types

// SecretPayload содержит данные для создания или получения секрета.
type SecretPayload struct {
	Name     string
	Type     SecretType
	Data     []byte
	Metadata *string
}
