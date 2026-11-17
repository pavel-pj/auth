package auth

import (
	"github.com/pavel-pj/auth/internal/service"
	desc "github.com/pavel-pj/auth/pkg/auth_v1"
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
