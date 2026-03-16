package service

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
)

//go:generate minimock -i Repository -o ./mocks/repo.go -n RepoMock
type Repository interface {
	WithTx(ctx context.Context, fn func(txRepo Repository) error) error
	// Close закрывает соединение с хранилищем и освобождает ресурсы
	Close() error

	CreateUser(ctx context.Context, user *types.User) (userID int64, err error)
	GetUserByLogin(ctx context.Context, login string) (user *types.User, err error)
	GetUserByID(ctx context.Context, id int64) (user *types.User, err error)

	GetSecrets(ctx context.Context, userID int64) (secrets []types.Secret, err error)
	CreateSecret(ctx context.Context, secret *types.Secret) error
	GetSecret(ctx context.Context, secretID, userID int64) (secret *types.Secret, err error)
	UpdateSecret(ctx context.Context, secret *types.Secret) error
	DeleteSecret(ctx context.Context, secretID, userID int64) error
}

type Option func(*Service)

func WithJWTSecret(jwtSecret []byte) Option {
	return func(s *Service) {
		s.jwtSecret = jwtSecret
	}
}

func WithMasterKey(masterKey []byte) Option {
	return func(s *Service) {
		s.masterKey = masterKey
	}
}

// Service реализует бизнес-логику работы с метриками.
type Service struct {
	repo      Repository
	logger    zerolog.Logger
	jwtSecret []byte
	masterKey []byte
}

func New(repo Repository, logger zerolog.Logger, opts ...Option) *Service {
	s := &Service{
		repo:   repo,
		logger: logger,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
