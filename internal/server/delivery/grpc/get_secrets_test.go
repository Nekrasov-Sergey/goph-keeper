package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	grpcServer "github.com/Nekrasov-Sergey/goph-keeper/internal/server/delivery/grpc"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/server/service"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/auth"
)

func TestServer_GetSecrets(t *testing.T) {
	t.Parallel()

	type want struct {
		secrets int
		err     bool
	}

	tests := []struct {
		name      string
		withToken bool
		buildMock func(*serviceMocks.RepoMock)
		want      want
	}{
		{
			name:      "Успешно",
			withToken: true,
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretsMock.Return([]types.Secret{
					{
						ID:   1,
						Name: "test",
						Type: types.SecretTypeText,
					},
				}, nil)
			},
			want: want{
				secrets: 1,
				err:     false,
			},
		},
		{
			name:      "Не авторизован",
			withToken: false,
			buildMock: func(repo *serviceMocks.RepoMock) {
			},
			want: want{
				err: true,
			},
		},
		{
			name:      "Неизвестный тип секрета",
			withToken: true,
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretsMock.Return([]types.Secret{
					{
						ID:   2,
						Name: "unknown",
						Type: types.SecretTypeUnknown,
					},
				}, nil)
			},
			want: want{
				err: true,
			},
		},
		{
			name:      "Внутренняя ошибка",
			withToken: true,
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretsMock.Return(nil, errors.New("some error"))
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)

			repo := serviceMocks.NewRepoMock(ctrl)
			if tt.buildMock != nil {
				tt.buildMock(repo)
			}

			logger := zerolog.Nop()
			svc := service.New(repo, logger)

			grpcSrv, err := grpcServer.New(svc, logger)
			require.NoError(t, err)

			lis, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

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
			require.NoError(t, err)

			t.Cleanup(func() {
				_ = conn.Close()
			})

			client := pb.NewKeeperClient(conn)

			ctx := context.Background()

			if tt.withToken {
				token, err := auth.GenerateToken(1, []byte(jwtSecret))
				require.NoError(t, err)

				md := metadata.New(map[string]string{
					"authorization": "Bearer " + token,
				})

				ctx = metadata.NewOutgoingContext(ctx, md)
			}

			resp, err := client.GetSecrets(ctx, &pb.GetSecretsRequest{})

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, resp.Secrets, tt.want.secrets)
		})
	}
}
