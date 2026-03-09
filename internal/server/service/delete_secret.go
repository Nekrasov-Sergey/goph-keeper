package service

import (
	"context"
)

func (s *Service) DeleteSecret(ctx context.Context, secretID, userID int64) error {
	if _, err := s.repo.GetSecret(ctx, secretID, userID); err != nil {
		return err
	}
	return s.repo.DeleteSecret(ctx, secretID, userID)
}
