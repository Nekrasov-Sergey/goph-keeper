package grpc_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	serviceMocks "github.com/Nekrasov-Sergey/goph-keeper/internal/server/service/mocks"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

func TestServer_Login(t *testing.T) {
	t.Parallel()

	type want struct {
		err bool
	}

	tests := []struct {
		name         string
		buildMock    func(*serviceMocks.RepoMock)
		loginRequest *pb.LoginRequest
		want         want
	}{
		{
			name: "Успешно",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetUserByLoginMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "$2a$10$myW3N2xys.qcVaQ5/DLnTOJbwkfgKrXcxgSu2qJHG/T7Q9m35nG5q",
					EncryptedUserKey: nil,
				}, nil)
			},
			loginRequest: &pb.LoginRequest{
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
			loginRequest: &pb.LoginRequest{
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
			loginRequest: &pb.LoginRequest{
				Login:    "test",
				Password: "",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Логин не существует",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetUserByLoginMock.Return(nil, errcodes.ErrInvalidCredentials)
			},
			loginRequest: &pb.LoginRequest{
				Login:    "wrong test",
				Password: "test",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Пароль неверный",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetUserByLoginMock.Return(&types.User{
					ID:               1,
					Login:            "test",
					Password:         "$2a$10$myW3N2xys.qcVaQ5/DLnTOJbwkfgKrXcxgSu2qJHG/T7Q9m35nG5q",
					EncryptedUserKey: nil,
				}, nil)
			},
			loginRequest: &pb.LoginRequest{
				Login:    "test",
				Password: "wrong test",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "Внутренняя ошибка",
			buildMock: func(repo *serviceMocks.RepoMock) {
				repo.GetUserByLoginMock.Return(nil, errors.New("some error"))
			},
			loginRequest: &pb.LoginRequest{
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

			resp, err := client.Login(context.Background(), tt.loginRequest)

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, resp.Token)
		})
	}
}
