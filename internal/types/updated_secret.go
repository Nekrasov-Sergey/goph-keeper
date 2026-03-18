// Package types содержит доменные типы данных приложения.
package types

// UpdatedSecret содержит данные для обновления секрета.
type UpdatedSecret struct {
	ID       int64
	Name     *string
	Data     []byte
	Metadata *string
}
