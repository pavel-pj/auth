package service

import (
	"auth/internal/model"
	"context"
)

type UserService interface {
	Store(ctx context.Context, user model.User) (int64, error)
}
