package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

func (s *Service) Login(ctx context.Context, user *types.User) (token string, err error) {
	dbUser, err := s.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return "", errcodes.ErrInvalidCredentials
	}

	return s.generateToken(dbUser.ID)
}
