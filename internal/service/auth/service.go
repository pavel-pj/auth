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
	JWTSecret  []byte
}

func NewAuthService(repository repository.AuthRepository, jwtSecret []byte) service.AuthService {
	return &AuthService{
		Repository: repository,
		JWTSecret:  jwtSecret,
	}
}
