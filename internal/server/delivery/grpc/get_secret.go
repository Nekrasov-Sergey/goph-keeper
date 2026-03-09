package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
)

func (s *Server) GetSecret(ctx context.Context, in *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	secret, err := s.service.GetSecret(ctx, in.Id, userID)
	if err != nil {
		if errors.Is(err, errcodes.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	secretType, err := mappers.DomainSecretTypeToProto(secret.Type)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetSecretResponse{
		Name:     secret.Name,
		Type:     secretType,
		Data:     secret.Data,
		Metadata: secret.Metadata,
	}, nil
}
