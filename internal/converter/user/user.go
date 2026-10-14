package user

import (
	"auth/internal/model"
	desc "auth/pkg/user_v1"
)

func ToUserFromProto(data *desc.UserInfo) model.User {

	return model.User{
		Name:            data.Name,
		Email:           data.Email,
		Password:        data.Password,
		PasswordConfirm: data.PasswordConfirm,
		Roles:           data.Roles.String(),
	}

}
