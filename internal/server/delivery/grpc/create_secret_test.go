package grpc_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/crypto"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

func TestServer_CreateSecret(t *testing.T) {
	t.Parallel()

	type want struct {
		err bool
	}

	tests := []struct {
		name      string
		buildMock func(*serviceMocks.RepoMock)
		req       *pb.CreateSecretRequest
		want      want
	}{
		{
			name: "Успешно",
			buildMock: func(repo *serviceMocks.RepoMock) {
				encryptedUserKey, err := crypto.Encrypt([]byte(masterKey), []byte(userKey))
				require.NoError(t, err)

				repo.GetUserByIDMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "test",
					EncryptedUserKey: encryptedUserKey,
				}, nil)

				repo.CreateSecretMock.Return(nil)
			},
			req: &pb.CreateSecretRequest{
				Name:     "test",
				Type:     pb.SecretType_Text,
				Data:     []byte("test"),
				Metadata: nil,
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Отсутствует имя секрета",
			buildMock: func(repo *serviceMocks.RepoMock) {
			},
			req: &pb.CreateSecretRequest{
				Name:     "",
				Type:     pb.SecretType_Text,
				Data:     []byte("test"),
				Metadata: nil,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Неизвестный тип секрета",
			buildMock: func(repo *serviceMocks.RepoMock) {
			},
			req: &pb.CreateSecretRequest{
				Name:     "test",
				Type:     pb.SecretType_Unspecified,
				Data:     []byte("test"),
				Metadata: nil,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Имя секрета уже занято",
			buildMock: func(repo *serviceMocks.RepoMock) {
				encryptedUserKey, err := crypto.Encrypt([]byte(masterKey), []byte(userKey))
				require.NoError(t, err)

				repo.GetUserByIDMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "test",
					EncryptedUserKey: encryptedUserKey,
				}, nil)

				repo.CreateSecretMock.Return(errcodes.ErrSecretNameAlreadyExists)
			},
			req: &pb.CreateSecretRequest{
				Name:     "test",
				Type:     pb.SecretType_Text,
				Data:     []byte("test"),
				Metadata: nil,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Внутренняя ошибка",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetUserByIDMock.Return(nil, errors.New("some error"))
			},
			req: &pb.CreateSecretRequest{
				Name:     "test",
				Type:     pb.SecretType_Text,
				Data:     []byte("test"),
				Metadata: nil,
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepoMock(t)
			if tt.buildMock != nil {
				tt.buildMock(repo)
			}

			client := NewTestClient(t, repo)
			ctx := WithToken(t, 1)

			_, err := client.CreateSecret(ctx, tt.req)

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
