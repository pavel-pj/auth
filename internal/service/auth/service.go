package auth

import (
	"auth/internal/repository"
	"auth/internal/service"
)

type AuthService struct {
	Repository repository.AuthRepository
}

func NewAuthService(repository repository.AuthRepository) service.AuthService {
	return AuthService{
		Repository: repository,
	}
}
