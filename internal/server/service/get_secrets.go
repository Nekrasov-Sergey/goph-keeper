package service

import (
	"context"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
)

func (s *Service) GetSecrets(ctx context.Context, userID int64) (secrets []types.Secret, err error) {
	return s.repo.GetSecrets(ctx, userID)
}
