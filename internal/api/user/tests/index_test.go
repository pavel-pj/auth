package tests

import (
	userApi "auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	"context"
	"errors"
	"testing"
	"time"

	serviceMocks "auth/internal/service/mocks"
	desc "auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req desc.IndexRequest
	}

	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		offset = uint64(1)
		limit  = uint64(10)

		errTest = errors.New("TestError")

		req = desc.IndexRequest{
			Offset: offset,
			Limit:  limit,
		}
	)
	const count = 10
	protoUsers := make([]*desc.User, 0, count)
	modelUsers := make([]model.User, 0, count)

	for x := 0; x < count; x++ {

		id := gofakeit.Int64()
		name := gofakeit.Name()
		email := gofakeit.Email()
		role := "ADMIN"
		createdAt := time.Now()
		updatedAt := time.Now()

		protoUsers = append(protoUsers, &desc.User{
			Id: id,
			Info: &desc.UserData{
				Name:  name,
				Email: email,
				Roles: desc.Roles_ADMIN,
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})

		modelUsers = append(modelUsers, model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Roles:     role,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})

	}

	wants := &desc.IndexResponse{Users: protoUsers}

	tests := []struct {
		name            string
		args            args
		want            *desc.IndexResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success",
			args: args{ctx, req},
			want: wants,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.IndexMock.Expect(ctx, limit, offset).Return(modelUsers, nil)
				return mock
			},
		},

		{
			name: "error",
			args: args{ctx, req},
			want: &desc.IndexResponse{},
			err:  errTest,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.IndexMock.Expect(ctx, limit, offset).Return([]model.User{}, errTest)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			api := userApi.NewUserImplimentation(tt.userServiceMock(mc))
			resp, err := api.Index(tt.args.ctx, &tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, resp)

		})

	}

}
