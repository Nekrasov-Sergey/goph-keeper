package grpc_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
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

			repo := NewRepoMock(t)
			if tt.buildMock != nil {
				tt.buildMock(repo)
			}

			client := NewTestClient(t, repo)

			ctx := context.Background()
			if tt.withToken {
				ctx = WithToken(t, 1)
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
