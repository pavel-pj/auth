package repository

import (
	model "auth/internal/repository/user/model"
	"context"
)

type UserRepository interface {
	Store(ctx context.Context, users model.User) (int64, error)
	//Delete(ctx context.Context, id int64) error
	//Get(ctx context.Context, id int64) (*desc.User, error)
}
