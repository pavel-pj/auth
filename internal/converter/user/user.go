package user

import (
	"auth/internal/model"
	desc "auth/pkg/user_v1"
)

func ToUserFromProto(data *desc.UserInfo) model.User {

	return model.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Roles:    data.Roles.String(),
	}

}

func ToProtoUserFromService(data model.User) *desc.User {

	return &desc.User{
		Id: data.ID,
		Info: &desc.UserInfo{
			Name:  data.Name,
			Email: data.Email,
			Roles: ToProtoRoles(data.Roles),
		},
	}

}

// model → proto
func ToProtoRoles(role string) desc.Roles {
	switch role {
	case "ADMIN":
		return desc.Roles_ADMIN
	default:
		return desc.Roles_USER
	}
}

// proto → model
func ToModelRoles(role desc.Roles) string {
	switch role {
	case desc.Roles_ADMIN:
		return "ADMIN"
	default:
		return "USER"
	}
}
