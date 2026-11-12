package access

import (
	"auth/internal/service"
)

var (
	authprefix = "Bearer "
)

type AccessService struct {
	JWT service.JWT
}

func NewAccessService(jwt service.JWT) service.AccessService {
	return &AccessService{
		JWT: jwt,
	}
}
