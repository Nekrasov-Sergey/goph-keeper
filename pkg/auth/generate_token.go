// Package auth содержит функции для работы с JWT-токенами.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// GenerateToken создаёт JWT-токен для пользователя с сроком действия 24 часа.
func GenerateToken(userID int64, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, "не удалось создать jwt token")
	}

	return token, nil
}
