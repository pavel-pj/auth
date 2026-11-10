package user

import (
	"auth/internal/service"
	desc "auth/pkg/user_v1"
	"errors"
)

var (
	ErrorPasswwordConfirmation = errors.New("passwords and password confirmation are not equal")
)

type UserImplementation struct {
	desc.UnimplementedUserV1Server
	serv service.UserService
}

func NewUserImplimentation(userService service.UserService) *UserImplementation {
	return &UserImplementation{
		serv: userService,
	}
}
