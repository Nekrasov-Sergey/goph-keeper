// Package service реализует бизнес-логику приложения.
package service

import (
	"context"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
)

// GetSecrets возвращает список всех секретов пользователя.
func (s *Service) GetSecrets(ctx context.Context, userID int64) (secrets []types.Secret, err error) {
	return s.repo.GetSecrets(ctx, userID)
}
