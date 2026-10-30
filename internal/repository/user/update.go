package user

import (
	"auth/internal/model"
	"context"
)

func (r *UserRepository) Update(ctx context.Context, user model.UserPatch) (model.User, error) {
	return model.User{}, nil

}
