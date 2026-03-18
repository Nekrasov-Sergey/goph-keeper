package grpc_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

func TestServer_Register(t *testing.T) {
	t.Parallel()

	type want struct {
		err bool
	}

	tests := []struct {
		name            string
		buildMock       func(*serviceMocks.RepoMock)
		registerRequest *pb.RegisterRequest
		want            want
	}{
		{
			name: "Успешно",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.CreateUserMock.Return(1, nil)
			},
			registerRequest: &pb.RegisterRequest{
				Login:    "test",
				Password: "test",
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Отсутствует логин",
			buildMock: func(repo *serviceMocks.RepoMock) {
			},
			registerRequest: &pb.RegisterRequest{
				Login:    "",
				Password: "test",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Отсутствует пароль",
			buildMock: func(repo *serviceMocks.RepoMock) {
			},
			registerRequest: &pb.RegisterRequest{
				Login:    "test",
				Password: "",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Логин уже занят",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.CreateUserMock.Return(0, errcodes.ErrLoginAlreadyExists)
			},
			registerRequest: &pb.RegisterRequest{
				Login:    "test",
				Password: "test",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Внутренняя ошибка",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.CreateUserMock.Return(0, errors.New("some error"))
			},
			registerRequest: &pb.RegisterRequest{
				Login:    "test",
				Password: "test",
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

			resp, err := client.Register(context.Background(), tt.registerRequest)

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, resp.Token)
		})
	}
}
