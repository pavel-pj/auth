package converter

import (
	"auth/internal/model"
	modelRepo "auth/internal/repository/user/model"
)

/*
func ToRepoUserFromService(data model.User) modelRepo.User {

	return modelRepo.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Roles:    data.Roles,
	}

}
*/

func ToServiceFromRepoUser(data modelRepo.User) model.User {
	return model.User{
		ID:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		Roles:     data.Roles,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
