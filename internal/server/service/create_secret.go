// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
)

// CreateSecret создаёт новый секрет для пользователя.
func (s *Service) CreateSecret(ctx context.Context, secretPayload *types.SecretPayload, userID int64) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	userKey, err := crypto.Decrypt(s.masterKey, user.EncryptedUserKey)
	if err != nil {
		return errors.Wrap(err, "не удалось расшифровать ключ пользователя")
	}

	encryptedData, err := crypto.Encrypt(userKey, secretPayload.Data)
	if err != nil {
		return errors.Wrap(err, "не удалось зашифровать секрет")
	}

	secret := &types.Secret{
		Name:          secretPayload.Name,
		Type:          secretPayload.Type,
		EncryptedData: encryptedData,
		Metadata:      secretPayload.Metadata,
		UserID:        userID,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateSecret(ctx, secret)
}
