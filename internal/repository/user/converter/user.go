package converter

import (
	modelService "auth/internal/model"
	model "auth/internal/repository/user/model"
)

func ToRepoUserFromService(data modelService.User) model.User {

	return model.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Roles:    data.Roles,
	}

}
