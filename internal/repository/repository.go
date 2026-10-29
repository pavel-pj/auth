package repository

import (
	logModel "auth/internal/repository/log/model"
	userModel "auth/internal/repository/user/model"
	"context"
)

type UserRepository interface {
	Store(ctx context.Context, users userModel.User) (int64, error)
	Get(ctx context.Context, id int64) (userModel.User, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Update(ctx context.Context, id int64) (userModel.User, error)
}

type LogRepository interface {
	Store(ctx context.Context, logs logModel.Log) error
}
