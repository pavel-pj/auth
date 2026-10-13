package user

import (
	repo "auth/internal/repository"
	serv "auth/internal/service"
	desc "auth/pkg/user_v1"
	"context"
)

type UserService struct {
	Repository repo.UserRepository
}

func NewUserService(repository repo.UserRepository) serv.UserService {
	return UserService{
		Repository: repository,
	}
}

func (s UserService) Store(ctx context.Context, users *desc.UserInfo) (int64, error) {
	return 252, nil
}
