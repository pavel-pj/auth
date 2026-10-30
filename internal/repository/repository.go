package repository

import (
	"auth/internal/model"
	logModel "auth/internal/repository/log/model"
	"context"
)

type UserRepository interface {
	Index(ctx context.Context, limit, offset uint64) ([]model.User, error)
	Store(ctx context.Context, users model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (model.User, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Update(ctx context.Context, user model.UserPatch) (model.User, error)
}

type LogRepository interface {
	Store(ctx context.Context, logs logModel.Log) error
}
