package user

import (
	model "auth/internal/repository/user/model"
	"context"
)

func (r *UserRepository) Update(ctx context.Context, id int64) (model.User, error) {
	return model.User{}, nil
}
