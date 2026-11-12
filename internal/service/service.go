package service

import (
	"auth/internal/model"
	"context"
	"time"
)

var (
	RefreshTokenExpiration = 60 * time.Minute
	AccessTokenExpiration  = 5 * time.Minute
)

type JWT struct {
	RefreshTokenSecret []byte
	AccessTokenSecret  []byte
}

type UserService interface {
	Index(ctx context.Context, limit, offset uint64) ([]model.User, error)
	Store(ctx context.Context, user model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user model.UserPatch) (model.User, error)
}

type LogService interface {
	Store(ctx context.Context, log model.Log) error
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, handler_path string) error
}
