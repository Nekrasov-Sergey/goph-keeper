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

// UpdateSecret обрабатывает запрос на обновление секрета.
func (s *Server) UpdateSecret(ctx context.Context, in *pb.UpdateSecretRequest) (*pb.UpdateSecretResponse, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	updatedSecret := &types.UpdatedSecret{
		ID:       in.Id,
		Name:     in.Name,
		Data:     in.Data,
		Metadata: in.Metadata,
	}

	if err := s.service.UpdateSecret(ctx, updatedSecret, userID); err != nil {
		if errors.Is(err, errcodes.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateSecretResponse{}, nil
}
