package user

import (
	repo "auth/internal/repository"
	serv "auth/internal/service"
)

type UserService struct {
	Repository repo.UserRepository
}

func NewUserService(repository repo.UserRepository) serv.UserService {
	return UserService{
		Repository: repository,
	}
}
