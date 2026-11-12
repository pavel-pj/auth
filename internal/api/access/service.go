package auth

import (
	"auth/internal/service"
	desc "auth/pkg/access_v1"
)

type AccessImplementation struct {
	desc.UnimplementedAccessV1Server
	serv service.AccessService
}

func NewAccessImplimentation(accessService service.AccessService) *AccessImplementation {
	return &AccessImplementation{
		serv: accessService,
	}
}
