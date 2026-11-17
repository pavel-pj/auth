package user

import (
	"errors"
	"github.com/pavel-pj/auth/internal/service"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
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
