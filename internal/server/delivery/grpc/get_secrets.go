// Package grpc реализует gRPC-сервер для обработки запросов.
package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/mappers"
)

// GetSecrets обрабатывает запрос на получение списка секретов.
func (s *Server) GetSecrets(ctx context.Context, _ *pb.GetSecretsRequest) (*pb.GetSecretsResponse, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	secrets, err := s.service.GetSecrets(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbSecrets := make([]*pb.Secret, 0, len(secrets))
	for _, secret := range secrets {
		secretType, err := mappers.DomainSecretTypeToProto(secret.Type)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		pbSecrets = append(pbSecrets, &pb.Secret{
			Id:   secret.ID,
			Name: secret.Name,
			Type: secretType,
		})
	}

	return &pb.GetSecretsResponse{
		Secrets: pbSecrets,
	}, nil
}
