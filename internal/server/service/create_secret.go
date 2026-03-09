package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
)

func (s *Service) CreateSecret(ctx context.Context, secretInput *types.SecretInput, userID int64) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	userKey, err := crypto.Decrypt(s.masterKey, user.EncryptedUserKey)
	if err != nil {
		return errors.Wrap(err, "не удалось расшифровать ключ пользователя")
	}

	encryptedData, err := crypto.Encrypt(userKey, secretInput.Data)
	if err != nil {
		return errors.Wrap(err, "не удалось зашифровать секрет")
	}

	secret := &types.Secret{
		Name:          secretInput.Name,
		Type:          secretInput.Type,
		EncryptedData: encryptedData,
		Metadata:      secretInput.Metadata,
		UserID:        userID,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateSecret(ctx, secret)
}
