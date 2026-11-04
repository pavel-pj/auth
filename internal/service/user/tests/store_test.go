package tests

import (
	dbMocks "auth/internal/client/mocks"
	"auth/internal/model"
	"auth/internal/repository"
	repositoryMocks "auth/internal/repository/mocks"
	"auth/internal/service"
	serviceMocks "auth/internal/service/mocks"
	"errors"

	"github.com/pavel-pj/platform_lib/pkg/client/db"

	userService "auth/internal/service/user"
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type logServiceyMockFunc func(mc *minimock.Controller) service.LogService
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		data model.UserCreate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		password = gofakeit.Password(true, true, true, true, true, 16)
		email    = gofakeit.Email()
		roles    = "ADMIN"
	)

	user := model.UserCreate{
		Name:     name,
		Password: password,
		Email:    email,
		Roles:    roles,
	}

	userLog := model.Log{
		Level:   model.LevelInfo,
		Message: fmt.Sprintf("user_store.name:%s,email:%s.role:%s", name, email, roles),
	}

	repoError := errors.New("Repository error")
	logError := errors.New("Log error")

	var inTx bool

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		logServiceyMock    logServiceyMockFunc
		txmanagerMock      txManagerMockFunc
	}{
		{
			name: "success",
			args: args{
				ctx:  ctx,
				data: user,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(m *minimock.Controller) repository.UserRepository {

				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.StoreMock.Set(func(ctx context.Context, got model.UserCreate) (i1 int64, err error) {
					require.True(t, inTx, "Repository.Store вызван ВНЕ транзакции")
					require.Equal(t, user, got)
					return id, nil
				})
				//mock.StoreMock.Expect(ctx, user).Return(id, nil)
				return mock
			},
			logServiceyMock: func(m *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				mock.StoreMock.Expect(ctx, userLog).Return(nil)
				return mock
			},
			txmanagerMock: func(m *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(m)
				mock.ReadCommittedMock.Set(func(ctx context.Context, h db.Handler) error {
					inTx = true
					defer func() { inTx = false }()
					return h(ctx)
				})
				return mock
			},
		},

		{
			name: "Error Repository",
			args: args{
				ctx:  ctx,
				data: user,
			},
			want: 0,
			err:  repoError,
			userRepositoryMock: func(m *minimock.Controller) repository.UserRepository {

				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.StoreMock.Set(func(ctx context.Context, got model.UserCreate) (i1 int64, err error) {
					require.True(t, inTx, "Repository.Store вызван ВНЕ транзакции")
					require.Equal(t, user, got)
					return 0, repoError
				})
				return mock
			},
			logServiceyMock: func(m *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				return mock
			},
			txmanagerMock: func(m *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(m)
				mock.ReadCommittedMock.Set(func(ctx context.Context, h db.Handler) error {
					inTx = true
					defer func() { inTx = false }()
					return h(ctx)
				})
				return mock
			},
		},

		{
			name: "error Log",
			args: args{
				ctx:  ctx,
				data: user,
			},
			want: 0,
			err:  logError,
			userRepositoryMock: func(m *minimock.Controller) repository.UserRepository {

				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.StoreMock.Set(func(ctx context.Context, got model.UserCreate) (i1 int64, err error) {
					require.True(t, inTx, "Repository.Store вызван ВНЕ транзакции")
					require.Equal(t, user, got)
					return id, nil
				})
				return mock
			},
			logServiceyMock: func(m *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				mock.StoreMock.Expect(ctx, userLog).Return(logError)
				return mock
			},
			txmanagerMock: func(m *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(m)
				mock.ReadCommittedMock.Set(func(ctx context.Context, h db.Handler) error {
					inTx = true
					defer func() { inTx = false }()
					return h(ctx)
				})
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

			resp, err := service.Store(tt.args.ctx, tt.args.data)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)

		})

	}

}
