// Package grpc реализует gRPC-сервер для обработки запросов.
package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

// Register обрабатывает запрос на регистрацию пользователя.
func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := &types.User{
		Login:    in.Login,
		Password: in.Password,
	}

	if err := validateUser(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.service.Register(ctx, user)
	if err != nil {
		if errors.Is(err, errcodes.ErrLoginAlreadyExists) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, err
	}

	return &pb.RegisterResponse{
		Token: token,
	}, nil
}

// validateUser проверяет корректность данных пользователя.
func validateUser(user *types.User) error {
	if user.Login == "" {
		return errors.New("отсутствует логин")
	}
	if user.Password == "" {
		return errors.New("отсутствует пароль")
	}
	return nil
}
