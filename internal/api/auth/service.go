package auth

import (
	"auth/internal/service"
	desc "auth/pkg/auth_v1"
)

type AuthImplementation struct {
	desc.UnimplementedAuthV1Server
	serv service.AuthService
}

func NewAuthImplimentation(authService service.AuthService) *AuthImplementation {
	return &AuthImplementation{
		serv: authService,
	}
}
