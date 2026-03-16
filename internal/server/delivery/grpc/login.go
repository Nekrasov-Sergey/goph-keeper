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

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := &types.User{
		Login:    in.Login,
		Password: in.Password,
	}

	if err := validateUser(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.service.Login(ctx, user)
	if err != nil {
		if errors.Is(err, errcodes.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}
