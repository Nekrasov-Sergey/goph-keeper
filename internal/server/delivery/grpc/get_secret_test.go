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
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

const (
	masterKey = "testtesttesttest"
	userKey   = "testtesttesttest"
	jwtSecret = "testtesttesttest"
)

func TestServer_GetSecret(t *testing.T) {
	t.Parallel()

	type want struct {
		secretName     string
		secretType     pb.SecretType
		secretData     []byte
		secretMetadata *string
		err            bool
	}

	tests := []struct {
		name      string
		buildMock func(*serviceMocks.RepoMock)
		want      want
	}{
		{
			name: "Успешно",
			buildMock: func(repo *serviceMocks.RepoMock) {
				encryptedUserKey, err := crypto.Encrypt([]byte(masterKey), []byte(userKey))
				require.NoError(t, err)

				encryptedData, err := crypto.Encrypt([]byte(userKey), []byte("some text"))
				require.NoError(t, err)

				repo.GetSecretMock.Return(&types.Secret{
					ID:            1,
					Name:          "test",
					Type:          types.SecretTypeText,
					EncryptedData: encryptedData,
					Metadata:      utils.Ptr("some metadata"),
					UserID:        1,
				}, nil)

				repo.GetUserByIDMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "test",
					EncryptedUserKey: encryptedUserKey,
				}, nil)
			},
			want: want{
				secretName:     "test",
				secretType:     pb.SecretType_Text,
				secretData:     []byte("some text"),
				secretMetadata: utils.Ptr("some metadata"),
				err:            false,
			},
		},
		{
			name: "Неизвестный тип секрета",
			buildMock: func(repo *serviceMocks.RepoMock) {
				encryptedUserKey, err := crypto.Encrypt([]byte(masterKey), []byte(userKey))
				require.NoError(t, err)

				encryptedData, err := crypto.Encrypt([]byte(userKey), []byte("some text"))
				require.NoError(t, err)

				repo.GetSecretMock.Return(&types.Secret{
					ID:            1,
					Name:          "test",
					Type:          types.SecretTypeUnknown,
					EncryptedData: encryptedData,
					Metadata:      utils.Ptr("some metadata"),
					UserID:        1,
				}, nil)

				repo.GetUserByIDMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "test",
					EncryptedUserKey: encryptedUserKey,
				}, nil)
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Секрет не найден",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(nil, errcodes.ErrSecretNotFound)
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Внутренняя ошибка",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(nil, errors.New("some error"))
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
			svc := service.New(repo, logger, service.WithMasterKey([]byte("testtesttesttest")))

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

			token, err := auth.GenerateToken(1, []byte(jwtSecret))
			require.NoError(t, err)

			md := metadata.New(map[string]string{
				"authorization": "Bearer " + token,
			})

			ctx := metadata.NewOutgoingContext(context.Background(), md)

			resp, err := client.GetSecret(ctx, &pb.GetSecretRequest{})

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.secretName, resp.Name)
			require.Equal(t, tt.want.secretType, resp.Type)
			require.Equal(t, tt.want.secretData, resp.Data)
			require.Equal(t, tt.want.secretMetadata, resp.Metadata)
		})
	}
}
