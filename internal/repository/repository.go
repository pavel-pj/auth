package repository

import (
	desc "auth/pkg/user_v1"
	"context"
)

type UserRepository interface {
	Store(ctx context.Context, users *desc.UserInfo) (int64, error)
	//Delete(ctx context.Context, id int64) error
	//Get(ctx context.Context, id int64) (*desc.User, error)
}
