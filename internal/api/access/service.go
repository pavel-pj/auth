package auth

import (
	"github.com/pavel-pj/auth/internal/service"
	desc "github.com/pavel-pj/auth/pkg/access_v1"
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
