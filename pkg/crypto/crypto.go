// Package crypto содержит функции шифрования и расшифрования данных.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/pkg/errors"
)

// Encrypt шифрует данные с использованием AES-GCM.
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать aes cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать gcm")
	}

	nonce, err := GenerateRandom(gcm.NonceSize())
	if err != nil {
		return nil, errors.Wrap(err, "не удалось сгенерировать nonce")
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt расшифровывает данные, зашифрованные с помощью Encrypt.
func Decrypt(key []byte, data []byte) ([]byte, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать aes cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать gcm")
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return nil, errors.New("некорректные зашифрованные данные")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось расшифровать данные")
	}

	return plaintext, nil
}

// GenerateRandom генерирует случайные байты заданного размера.
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)

	if _, err := rand.Read(b); err != nil {
		return nil, errors.Wrap(err, "не удалось сгенерировать случайные байты")
	}

	return b, nil
}

// validateKey проверяет, что ключ имеет допустимую длину (16, 24 или 32 байта).
func validateKey(key []byte) error {
	switch len(key) {
	case 16, 24, 32:
		return nil
	default:
		return errors.New("ключ должен быть 16, 24 или 32 байта")
	}
}
