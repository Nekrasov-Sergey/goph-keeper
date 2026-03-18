// Package service реализует бизнес-логику приложения.
package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/auth"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
)

// Register регистрирует нового пользователя и возвращает JWT-токен.
func (s *Service) Register(ctx context.Context, user *types.User) (token string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	userKey, err := crypto.GenerateRandom(32)
	if err != nil {
		return "", err
	}

	encryptedUserKey, err := crypto.Encrypt(s.masterKey, userKey)
	if err != nil {
		return "", err
	}
	user.EncryptedUserKey = encryptedUserKey

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return auth.GenerateToken(userID, s.jwtSecret)
}
