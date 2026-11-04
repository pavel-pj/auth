package tests

import (
	userApi "auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	serviceMocks "auth/internal/service/mocks"
	desc "auth/pkg/user_v1"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {

	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = "ADMIN"
		createdAt = time.Now()
		updatedAt = time.Now()

		userModel = model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Roles:     role,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		res = &desc.GetResponse{
			User: &desc.User{
				Id: id,
				Info: &desc.UserData{
					Name:  name,
					Email: email,
					Roles: desc.Roles_ADMIN,
				},
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}

		errorTest = errors.New("TestError")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: desc.GetRequest{
					Id: id,
				},
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(userModel, nil)
				return mock
			},
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				req: desc.GetRequest{
					Id: id,
				},
			},
			want: &desc.GetResponse{},
			err:  errorTest,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(model.User{}, errorTest)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()

			api := userApi.NewUserImplimentation(tt.userServiceMock(mc))
			resp, err := api.Get(tt.args.ctx, &tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, resp)

		})

	}

}
