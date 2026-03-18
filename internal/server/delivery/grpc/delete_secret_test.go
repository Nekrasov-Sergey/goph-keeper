package grpc_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

func TestServer_DeleteSecret(t *testing.T) {
	t.Parallel()

	type want struct {
		err bool
	}

	tests := []struct {
		name      string
		buildMock func(*serviceMocks.RepoMock)
		want      want
	}{
		{
			name: "Успешно",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetSecretMock.Return(nil, nil)
				repo.DeleteSecretMock.Return(nil)
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

			_, err := client.DeleteSecret(ctx, &pb.DeleteSecretRequest{
				Id: 1,
			})

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
