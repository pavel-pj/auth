package tests

import (
	"auth/internal/model"
	"auth/internal/repository"
	repositoryMocks "auth/internal/repository/mocks"
	"auth/internal/service"
	serviceMocks "auth/internal/service/mocks"
	userService "auth/internal/service/user"
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
)

func TestStore(t *testing.T) {

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type logServiceyMockFunc func(mc *minimock.Controller) service.LogService

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
		Message: fmt.Sprintf("user_store.name. some info"),
	}

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		logServiceyMock    logServiceyMockFunc
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
				mock.StoreMock.Expect(ctx, user).Return(id, nil)
				return mock
			},
			logServiceyMock: func(m *minimock.Controller) service.LogService {
				mock := serviceMocks.NewLogServiceMock(mc)
				mock.StoreMock.Expect(ctx, userLog).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			service := userService.NewUserService(tt.userRepositoryMock(mc), tt.logServiceyMock)

		})

	}

}
