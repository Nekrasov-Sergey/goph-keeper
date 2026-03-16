package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
)

func (s *Server) CreateSecret(ctx context.Context, in *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "отсутствует имя секрета")
	}

	secretType, err := mappers.ProtoSecretTypeToDomain(in.Type)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	secretPayload := &types.SecretPayload{
		Name:     in.Name,
		Type:     secretType,
		Data:     in.Data,
		Metadata: in.Metadata,
	}

	err = s.service.CreateSecret(ctx, secretPayload, userID)
	if err != nil {
		if errors.Is(err, errcodes.ErrSecretNameAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateSecretResponse{}, nil
}
