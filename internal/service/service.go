package service

import (
	desc "auth/pkg/user_v1"
	"context"
)

type UserService interface {
	Store(ctx context.Context, user *desc.UserInfo) (int64, error)
}
