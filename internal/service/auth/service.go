package auth

import (
	"errors"
	"github.com/pavel-pj/auth/internal/repository"
	"github.com/pavel-pj/auth/internal/service"
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
