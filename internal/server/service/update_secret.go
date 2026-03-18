// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

// UpdateSecret обновляет секрет пользователя.
func (s *Service) UpdateSecret(ctx context.Context, updatedSecret *types.UpdatedSecret, userID int64) error {
	currentSecret, err := s.repo.GetSecret(ctx, updatedSecret.ID, userID)
	if err != nil {
		return err
	}

	if updatedSecret.Name != nil {
		currentSecret.Name = utils.Deref(updatedSecret.Name)
	}

	if updatedSecret.Data != nil {
		user, err := s.repo.GetUserByID(ctx, userID)
		if err != nil {
			return err
		}

		userKey, err := crypto.Decrypt(s.masterKey, user.EncryptedUserKey)
		if err != nil {
			return errors.Wrap(err, "не удалось расшифровать ключ пользователя")
		}

		currentSecret.EncryptedData, err = crypto.Encrypt(userKey, updatedSecret.Data)
		if err != nil {
			return errors.Wrap(err, "не удалось зашифровать данные")
		}
	}

	if updatedSecret.Metadata != nil {
		currentSecret.Metadata = updatedSecret.Metadata
	}

	secret := &types.Secret{
		ID:            updatedSecret.ID,
		Name:          currentSecret.Name,
		Type:          currentSecret.Type,
		EncryptedData: currentSecret.EncryptedData,
		Metadata:      currentSecret.Metadata,
		UpdatedAt:     utils.Ptr(time.Now()),
	}

	return s.repo.UpdateSecret(ctx, secret)
}
