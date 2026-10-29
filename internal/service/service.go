package service

import (
	"auth/internal/model"
	"context"
	"errors"
)

var ErrNotFound = errors.New("entity not found")

type UserService interface {
	Store(ctx context.Context, user model.User) (int64, error)
	Get(ctx context.Context, id int64) (model.User, error)
	Delete(ctx context.Context, id int64) error
}

type LogService interface {
	Store(ctx context.Context, log model.Log) error
}
