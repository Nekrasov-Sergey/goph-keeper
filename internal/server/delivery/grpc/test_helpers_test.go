// Package grpc_test содержит тесты для gRPC-сервера.
package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	grpcServer "github.com/Nekrasov-Sergey/goph-keeper/internal/server/delivery/grpc"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/server/service"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/auth"
)

// Тестовые константы.
const (
	masterKey = "testtesttesttesttesttesttesttest" // 32 байта для AES-256
	userKey   = "testtesttesttesttesttesttesttest" // 32 байта для AES-256
	jwtSecret = "testtesttesttesttesttesttesttest" // 32 байта
)

// NewTestClient создаёт тестовый gRPC-клиент с запущенным сервером.
func NewTestClient(t *testing.T, repo *serviceMocks.RepoMock) pb.KeeperClient {
	t.Helper()

	logger := zerolog.Nop()
	svc := service.New(repo, logger, service.WithMasterKey([]byte(masterKey)))

	grpcSrv, err := grpcServer.New(svc, logger)
	if err != nil {
		t.Fatalf("failed to create gRPC server: %v", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcServer.AuthInterceptor([]byte(jwtSecret))),
	)
	pb.RegisterKeeperServer(server, grpcSrv)

	go func() {
		_ = server.Serve(lis)
	}()
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return pb.NewKeeperClient(conn)
}

// WithToken добавляет JWT-токен в контекст запроса.
func WithToken(t *testing.T, userID int64) context.Context {
	t.Helper()

	token, err := auth.GenerateToken(userID, []byte(jwtSecret))
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	return metadata.NewOutgoingContext(context.Background(), md)
}

// NewRepoMock создаёт новый mock репозитория.
func NewRepoMock(t *testing.T) *serviceMocks.RepoMock {
	t.Helper()
	return serviceMocks.NewRepoMock(minimock.NewController(t))
}
