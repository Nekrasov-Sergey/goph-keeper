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
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

func TestServer_UpdateSecret(t *testing.T) {
	t.Parallel()

	type want struct {
		err bool
	}

	tests := []struct {
		name      string
		buildMock func(*serviceMocks.RepoMock)
		req       *pb.UpdateSecretRequest
		want      want
	}{
		{
			name: "Успешное обновление названия",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(&types.Secret{
					Name: "current name",
				}, nil)

				repo.UpdateSecretMock.Return(nil)
			},
			req: &pb.UpdateSecretRequest{
				Id:   1,
				Name: utils.Ptr("new name"),
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Успешное обновление данных",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(&types.Secret{
					EncryptedData: []byte("current encrypted data"),
				}, nil)

				encryptedUserKey, err := crypto.Encrypt([]byte(masterKey), []byte(userKey))
				require.NoError(t, err)

				repo.GetUserByIDMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "test",
					EncryptedUserKey: encryptedUserKey,
				}, nil)

				repo.UpdateSecretMock.Return(nil)
			},
			req: &pb.UpdateSecretRequest{
				Id:   1,
				Data: []byte("new data"),
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Успешное обновление метаданных",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(&types.Secret{
					Metadata: utils.Ptr("current metadata"),
				}, nil)

				repo.UpdateSecretMock.Return(nil)
			},
			req: &pb.UpdateSecretRequest{
				Id:       1,
				Metadata: utils.Ptr("new metadata"),
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Секрет не найден",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(nil, errcodes.ErrSecretNotFound)
			},
			req: &pb.UpdateSecretRequest{
				Id:   2,
				Name: utils.Ptr("new name"),
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
			req: &pb.UpdateSecretRequest{
				Id:   2,
				Name: utils.Ptr("new name"),
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

			_, err := client.UpdateSecret(ctx, tt.req)

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
