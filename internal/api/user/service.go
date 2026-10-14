package user

import (
	"auth/internal/service"
	desc "auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	serv service.UserService
}

func NewUserImplimentation(userService service.UserService) *Implementation {
	return &Implementation{
		serv: userService,
	}
}
