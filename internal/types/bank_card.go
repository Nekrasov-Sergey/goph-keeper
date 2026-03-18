// Package types содержит доменные типы данных приложения.
package types

// BankCard представляет данные банковской карты для хранения в секрете.
type BankCard struct {
	Number string `json:"number"`
	Holder string `json:"holder"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv"`
}
