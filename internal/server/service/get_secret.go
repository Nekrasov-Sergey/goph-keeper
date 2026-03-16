package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
)

func (s *Service) GetSecret(ctx context.Context, secretID, userID int64) (*types.SecretPayload, error) {
	secret, err := s.repo.GetSecret(ctx, secretID, userID)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userKey, err := crypto.Decrypt(s.masterKey, user.EncryptedUserKey)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось расшифровать ключ пользователя")
	}

	data, err := crypto.Decrypt(userKey, secret.EncryptedData)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось расшифровать секрет")
	}

	secretPayload := &types.SecretPayload{
		Name:     secret.Name,
		Type:     secret.Type,
		Data:     data,
		Metadata: secret.Metadata,
	}

	return secretPayload, nil
}
