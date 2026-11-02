package tests

import (
	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	desc "auth/pkg/user_v1"
	"fmt"

	"context"
	"testing"

	serviceMocks "auth/internal/service/mocks"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
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
		password = "Password04%"
		//passwordFail = "Password04_"
		role = "ADMIN"

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

		//===============================================
		/*
			reqFailConfirmPassword = &desc.StoreRequest{
				User: &desc.UserCreate{

					Info: &desc.UserData{
						Name:  name,
						Email: email,
						Roles: 1,
					},
					Password:        password,
					PasswordConfirm: passwordFail,
				},
			}
		*/
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
		/*
			{
				name: "error confirmation password",
				args: args{
					ctx: ctx,
					req: reqFailConfirmPassword,
				},
				want: &desc.StoreResponse{},
				err:  api.ErrorPasswwordConfirmation,
				userServiceMock: func(mc *minimock.Controller) service.UserService {

					return serviceMocks.NewUserServiceMock(mc)
				},
			},
		*/
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//userServiceMock := tt.userServiceMock
			api := user.NewUserImplimentation(tt.userServiceMock(mc))
			res, err := api.Store(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want.Id, res.Id)

		})

	}

}
