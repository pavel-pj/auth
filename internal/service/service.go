package service

import (
	"auth/internal/model"
	"context"
)

type UserService interface {
	Store(ctx context.Context, user model.User) (int64, error)
	Get(ctx context.Context, id int64) (model.User, error)
}

type LogService interface {
	Store(ctx context.Context, log model.Log) error
}
