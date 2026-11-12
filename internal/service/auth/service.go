package auth

import (
	"auth/internal/repository"
	"auth/internal/service"
	"errors"
)

var (
	errorWrongPassword = errors.New("Incorrect passwrod for user")
)

type AuthService struct {
	Repository repository.AuthRepository
	JWT        service.JWT
}

func NewAuthService(repository repository.AuthRepository, jwt service.JWT) service.AuthService {
	return &AuthService{
		Repository: repository,
		JWT:        jwt,
	}
}
