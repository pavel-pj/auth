package user

import (
	"context"
	"log"

	model "auth/internal/model"
	converter "auth/internal/repository/user/converter"

	"github.com/sppps/pwdhash"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s UserService) Store(ctx context.Context, user model.User) (int64, error) {

	if user.Password != user.PasswordConfirm {
		return 0, status.Error(codes.InvalidArgument, "password confirmation does not match")
	}

	hash, err := pwdhash.Hash(user.Password)
	if err != nil {
		log.Fatal(err)
		return 0, nil
	}

	hashedUser := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hash,
		Roles:    user.Roles,
	}

	return s.Repository.Store(ctx, converter.ToRepoUserFromService(hashedUser))

}
