// Package service реализует бизнес-логику приложения.
package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// generateToken создаёт JWT-токен для пользователя.
func (s *Service) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, "не удалось создать jwt token")
	}

	return token, nil
}
