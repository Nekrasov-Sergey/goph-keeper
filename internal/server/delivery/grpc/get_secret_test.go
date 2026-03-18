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

			repo := NewRepoMock(t)
			if tt.buildMock != nil {
				tt.buildMock(repo)
			}

			client := NewTestClient(t, repo)
			ctx := WithToken(t, 1)

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
