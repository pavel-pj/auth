package tests

import (
	dbMocks "auth/internal/client/mocks"
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
	"errors"
	"testing"
	"time"

	repositoryMocks "auth/internal/repository/mocks"
	serviceMocks "auth/internal/service/mocks"
	userService "auth/internal/service/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/pavel-pj/platform_lib/pkg/client/db"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {

	type userReposytoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type logServiceyMockFunc func(mc *minimock.Controller) service.LogService
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = int64(252)

		resp = model.User{
			ID:        id,
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Roles:     "ADMIN",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		repoError = errors.New("repository error ")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		want               model.User
		userRepositoryMock userReposytoryMockFunc
		logServiceyMock    logServiceyMockFunc
		txmanagerMock      txManagerMockFunc
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: id,
			},
			err:  nil,
			want: resp,
			userRepositoryMock: func(m *minimock.Controller) repository.UserRepository {

				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(resp, nil)
				return mock
			},
			logServiceyMock: func(m *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				return mock
			},
			txmanagerMock: func(m *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(m)
				return mock
			},
		},
		{
			name: "repository error",
			args: args{
				ctx: ctx,
				req: id,
			},
			err:  repoError,
			want: model.User{},
			userRepositoryMock: func(m *minimock.Controller) repository.UserRepository {

				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(model.User{}, repoError)
				return mock
			},
			logServiceyMock: func(m *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				return mock
			},
			txmanagerMock: func(m *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(m)
				return mock
			},
		},
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			service := userService.NewUserService(
				tt.userRepositoryMock(mc),
				tt.txmanagerMock(mc),
				tt.logServiceyMock(mc),
			)

			res, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})

	}

}
