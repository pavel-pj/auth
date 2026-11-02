package tests

import (
	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	desc "auth/pkg/user_v1"
	"context"
	"fmt"
	"testing"

	serviceMocks "auth/internal/service/mocks"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.StoreRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 16)
		role     = "ADMIN"

		serviceErr = fmt.Errorf("service error")

		req = &desc.StoreRequest{
			User: &desc.UserCreate{

				Info: &desc.UserData{
					Name:  name,
					Email: email,
					Roles: 1,
				},
				Password:        password,
				PasswordConfirm: password,
			},
		}

		info = model.UserCreate{
			Name:     name,
			Email:    email,
			Password: password,
			Roles:    role,
		}

		res = &desc.StoreResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.StoreResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.StoreMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},

		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &desc.StoreResponse{},
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.StoreMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//userServiceMock := tt.userServiceMock
			api := user.NewUserImplimentation(tt.userServiceMock(mc))
			res, err := api.Store(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want.Id, res.Id)

		})

	}

}
